package cvrender_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/cvrender"
	render_html "github.com/germainlefebvre4/cvwonder/internal/cvrender/html"
	"github.com/germainlefebvre4/cvwonder/internal/cvrender/outpath"
	"github.com/germainlefebvre4/cvwonder/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputPaths(t *testing.T) {
	assert.Equal(t, []string{filepath.Join("out", "cv.html")}, cvrender.OutputPaths("out", "dir/cv.yml", "html"))
	assert.Equal(t, []string{filepath.Join("out", "cv.html"), filepath.Join("out", "cv.pdf")}, cvrender.OutputPaths("out", "dir/cv.yml", "pdf"))
}

// The helper must name exactly the file the HTML renderer writes.
func TestOutputPathsMatchHTMLRenderer(t *testing.T) {
	base := t.TempDir()
	themeDir := filepath.Join(base, "themes", "demo")
	require.NoError(t, os.MkdirAll(themeDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(themeDir, "index.html"), []byte("<html>{{ .Person.Name }}</html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(themeDir, "theme.yaml"), []byte("name: Demo\nslug: demo\n"), 0o644))
	out := filepath.Join(base, "generated")
	require.NoError(t, os.MkdirAll(out, 0o755))

	svc, err := render_html.NewRenderHTMLServices()
	require.NoError(t, err)
	require.NoError(t, svc.RenderFormatHTML(model.CV{}, base, out, outpath.Base("some/dir/my-cv.yml"), "demo", false, nil))

	want := cvrender.OutputPaths(out, "some/dir/my-cv.yml", "html")
	require.Len(t, want, 1)
	assert.FileExists(t, want[0])
	html, err := filepath.Glob(filepath.Join(out, "*.html"))
	require.NoError(t, err)
	assert.Equal(t, want, html, "the renderer writes exactly one HTML file, the one the helper names")
}
