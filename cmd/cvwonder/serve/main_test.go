package cmdServe

import (
	"bytes"
	"errors"
	"testing"

	"github.com/charmbracelet/colorprofile"
	rendermocks "github.com/germainlefebvre4/cvwonder/internal/cvrender/mocks"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func reporter() (*ui.Text, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return ui.NewTextProfile(out, &bytes.Buffer{}, colorprofile.NoTTY), out
}

func TestReportingRenderer_SuccessAndFailure(t *testing.T) {
	rep, out := reporter()
	inner := rendermocks.NewRenderInterfaceMock(t)
	inner.EXPECT().Render(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	inner.EXPECT().Render(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("bad template")).Once()
	r := reportingRenderer{RenderInterface: inner, rep: rep, input: "cv.yml"}

	require.NoError(t, r.Render(model.CV{}, "", "", "", "", "html", true, nil))
	assert.Equal(t, "✓ Re-rendered cv.yml\n", out.String())

	out.Reset()
	require.Error(t, r.Render(model.CV{}, "", "", "", "", "html", true, nil))
	assert.Empty(t, out.String(), "failures are reported through the watcher error channel")
}

func TestReportWatchErrors_StepPerFailure(t *testing.T) {
	rep, out := reporter()
	errs := make(chan error, 2)
	errs <- errors.New("parsing cv.yml: bad yaml")
	errs <- errors.New("rendering cv.yml: boom")
	close(errs)
	reportWatchErrors(rep, errs) // returns once the channel is closed

	assert.Equal(t, "✗ Re-render failed: parsing cv.yml: bad yaml\n✗ Re-render failed: rendering cv.yml: boom\n", out.String())
}
