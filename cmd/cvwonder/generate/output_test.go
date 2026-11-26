package cmdGenerate

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

func newReporter() (*ui.Text, *bytes.Buffer, *bytes.Buffer) {
	out, errb := &bytes.Buffer{}, &bytes.Buffer{}
	return ui.NewTextProfile(out, errb, colorprofile.NoTTY), out, errb
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

	g := &generator{rep: rep, validator: val, parser: parser, renderer: renderer}
	require.NoError(t, g.generateSingle())

	got := out.String()
	assert.Empty(t, errb.String())
	assert.True(t, strings.HasPrefix(got, "CV Wonder\n"), got)
	order := []string{"Input:", "Output:", "Theme:", "Format:", "Validation passed", "Parsed cv.yml", "Rendered html", "CV generated successfully"}
	last := -1
	for _, want := range order {
		i := strings.Index(got, want)
		require.GreaterOrEqual(t, i, 0, "missing %q in %q", want, got)
		require.Greater(t, i, last, "%q out of order in %q", want, got)
		last = i
	}
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

	g := &generator{rep: rep, parser: parser, renderer: renderer}
	err := g.generateSingle()
	require.ErrorContains(t, err, "boom")
	assert.Contains(t, out.String(), "✗ Rendering html")
	assert.NotContains(t, out.String(), "CV generated successfully")
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
