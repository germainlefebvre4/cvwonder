package render_screenshot

import (
	"testing"

	serveMocks "github.com/germainlefebvre4/cvwonder/internal/cvserve/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewRenderScreenshotServices(t *testing.T) {
	t.Run("Should create RenderScreenshotServices successfully", func(t *testing.T) {
		serveMock := serveMocks.NewServeInterfaceMock(t)

		service, err := NewRenderScreenshotServices(serveMock)

		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}
