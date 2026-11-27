package cmdGenerate

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/charmbracelet/colorprofile"
	"github.com/germainlefebvre4/cvwonder/internal/cvbulk"
	parsermocks "github.com/germainlefebvre4/cvwonder/internal/cvparser/mocks"
	rendermocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/mocks"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/germainlefebvre4/cvwonder/internal/utils"
	"github.com/germainlefebvre4/cvwonder/internal/validator"
	validatormocks "github.com/germainlefebvre4/cvwonder/internal/validator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// stepClock advances 10ms on every call, so durations are deterministic.
func stepClock() func() time.Time {
	t := time.Unix(1000, 0)
	return func() time.Time { t = t.Add(10 * time.Millisecond); return t }
}

func newReporter() (*ui.Text, *bytes.Buffer, *bytes.Buffer) {
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	return ui.NewTextProfile(out, errb, colorprofile.NoTTY, ui.WithClock(stepClock()), ui.WithVersion("1.2.3")), out, errb
}

// singleFixture prepares a working directory with a theme and an input file and
// sets utils.CliArgs for a single-file run.
func singleFixture(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	themeDir := filepath.Join(dir, "themes", "default")
	require.NoError(t, os.MkdirAll(themeDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(themeDir, "theme.yaml"), []byte("name: Default\nslug: default\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "cv.yml"), []byte("person:\n  name: Jane\n"), 0o644))
	t.Chdir(dir)
	prev := utils.CliArgs
	utils.CliArgs = utils.Configuration{InputFile: "cv.yml", OutputDirectory: "generated/", ThemeName: "default", Format: "html"}
	t.Cleanup(func() { utils.CliArgs = prev })
}

func TestGenerateSingle_OrderedStepsOnStdout(t *testing.T) {
	singleFixture(t)
	utils.CliArgs.Validate = true
	rep, out, errb := newReporter()

	val := validatormocks.NewValidatorInterfaceMock(t)
	val.EXPECT().ValidateFile(mock.Anything).Return(&validator.ValidationResult{Valid: true}, nil)
	parser := parsermocks.NewParserInterfaceMock(t)
	parser.EXPECT().ParseFile(mock.Anything).Return(model.CV{}, nil)
	renderer := rendermocks.NewRenderInterfaceMock(t)
	renderer.EXPECT().Render(mock.Anything, mock.Anything, mock.Anything, mock.Anything, "default", "html", false, mock.Anything).Return(nil)

	g := &generator{rep: rep, validator: val, parser: parser, renderer: renderer, now: stepClock()}
	require.NoError(t, g.generateSingle())

	want := "CV Wonder v1.2.3\n" +
		"  › Input   cv.yml\n" +
		"  › Output  generated/\n" +
		"  › Theme   default\n" +
		"  › Format  html\n" +
		"\n" +
		"✓ Validation passed: the CV YAML file is valid\n" +
		"✓ Parsed cv.yml  10ms\n" +
		"✓ Rendered html  10ms\n" +
		"\n" +
		"Done in 10ms\n" +
		"  " + filepath.Join("generated", "cv.html") + "\n" +
		"\n" +
		"Next: cvwonder serve --watch\n"
	assert.Empty(t, errb.String())
	assert.Equal(t, want, out.String())
	assert.NotContains(t, out.String(), "CV generated successfully")
}

func TestGenerateSingle_PDFFooterListsHTMLAndPDF(t *testing.T) {
	singleFixture(t)
	utils.CliArgs.Format = "pdf"
	rep, out, _ := newReporter()
	parser := parsermocks.NewParserInterfaceMock(t)
	parser.EXPECT().ParseFile(mock.Anything).Return(model.CV{}, nil)
	renderer := rendermocks.NewRenderInterfaceMock(t)
	renderer.EXPECT().Render(mock.Anything, mock.Anything, mock.Anything, mock.Anything, "default", "pdf", false, mock.Anything).Return(nil)

	g := &generator{rep: rep, parser: parser, renderer: renderer, now: stepClock()}
	require.NoError(t, g.generateSingle())
	got := out.String()
	assert.Contains(t, got, "  "+filepath.Join("generated", "cv.html")+"\n  "+filepath.Join("generated", "cv.pdf")+"\n")
	assert.Contains(t, got, "Next: cvwonder serve --watch\n")
	assert.NotContains(t, got, "\x1b")
}

func TestGenerateSingle_NextStepRepeatsNonDefaultFlags(t *testing.T) {
	singleFixture(t)
	utils.CliArgs.OutputDirectory = "public/"
	utils.CliArgs.ThemeName = "default"
	assert.Equal(t, "cvwonder serve --watch -i my.yml -o public/", serveHint("my.yml", "public/"))
	assert.Equal(t, "cvwonder serve --watch", serveHint("cv.yml", "generated/"))
}

func TestGenerateSingle_QuietPrintsNoFooterOrHint(t *testing.T) {
	singleFixture(t)
	rep, out, _ := newReporter()
	parser := parsermocks.NewParserInterfaceMock(t)
	parser.EXPECT().ParseFile(mock.Anything).Return(model.CV{}, nil)
	renderer := rendermocks.NewRenderInterfaceMock(t)
	renderer.EXPECT().Render(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	g := &generator{rep: ui.Quiet(rep), parser: parser, renderer: renderer, now: stepClock()}
	require.NoError(t, g.generateSingle())
	assert.Empty(t, out.String())
}

func TestGenerateSingle_InvalidFileShowsFindingsAndStops(t *testing.T) {
	singleFixture(t)
	utils.CliArgs.Validate = true
	rep, out, _ := newReporter()

	val := validatormocks.NewValidatorInterfaceMock(t)
	val.EXPECT().ValidateFile(mock.Anything).Return(&validator.ValidationResult{
		Valid:  false,
		Errors: []validator.ValidationError{{Field: "person.name", Message: "required", LineNumber: 2}},
	}, nil)
	// No parser/renderer expectations: any call fails the test.
	g := &generator{rep: rep, validator: val, parser: parsermocks.NewParserInterfaceMock(t), renderer: rendermocks.NewRenderInterfaceMock(t)}

	err := g.generateSingle()
	require.ErrorIs(t, err, ErrValidationFailed)
	assert.Contains(t, out.String(), "✗ Validation failed: 1 error")
	assert.Contains(t, out.String(), "Field: person.name")
	assert.NoDirExists(t, "generated")
}

func TestGenerateSingle_RenderFailureIsAnError(t *testing.T) {
	singleFixture(t)
	rep, out, _ := newReporter()
	parser := parsermocks.NewParserInterfaceMock(t)
	parser.EXPECT().ParseFile(mock.Anything).Return(model.CV{}, nil)
	renderer := rendermocks.NewRenderInterfaceMock(t)
	renderer.EXPECT().Render(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("boom"))

	g := &generator{rep: rep, parser: parser, renderer: renderer, now: stepClock()}
	err := g.generateSingle()
	require.ErrorContains(t, err, "boom")
	assert.Contains(t, out.String(), "✗ Rendering html  10ms\n")
	assert.NotContains(t, out.String(), "Done in")
	assert.NotContains(t, out.String(), "Next:")
}

func TestReportBulk_Mixed(t *testing.T) {
	rep, out, _ := newReporter()
	reportBulk(rep, []cvbulk.BulkResult{
		{FilePath: "c.yml", Success: true}, // results arrive in worker completion order
		{FilePath: "b.yml", Success: false, Error: errors.New("parse error")},
		{FilePath: "a.yml", Success: true},
	})
	want := "✓ a.yml\n✗ b.yml -- parse error\n✓ c.yml\n! Total: 3 | Success: 2 | Failed: 1\n"
	assert.Equal(t, want, out.String())
}

func TestReportBulk_AllSuccess(t *testing.T) {
	rep, out, _ := newReporter()
	reportBulk(rep, []cvbulk.BulkResult{{FilePath: "a.yml", Success: true}, {FilePath: "b.yml", Success: true}})
	assert.Equal(t, "✓ a.yml\n✓ b.yml\n✓ Total: 2 | Success: 2 | Failed: 0\n", out.String())
}

func TestReportBulk_QuietStillPrintsReport(t *testing.T) {
	rep, out, _ := newReporter()
	q := ui.Quiet(rep)
	q.Header("CV Wonder - Bulk Mode")
	reportBulk(q, []cvbulk.BulkResult{{FilePath: "a.yml", Success: true}})
	assert.Equal(t, "✓ a.yml\n✓ Total: 1 | Success: 1 | Failed: 0\n", out.String())
}
