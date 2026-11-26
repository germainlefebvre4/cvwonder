package cmdThemes

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"

	"github.com/germainlefebvre4/cvwonder/internal/themes"
	"github.com/germainlefebvre4/cvwonder/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestInstallError_RefSwitchHintUnlessForced(t *testing.T) {
	base := fmt.Errorf("%w 'v2': %w", themes.ErrRefSwitch, errors.New("local changes"))

	assert.Equal(t, ui.ForceHint, ui.HintOf(installError(base, false)))
	assert.Empty(t, ui.HintOf(installError(base, true)))
	assert.Empty(t, ui.HintOf(installError(errors.New("network down"), false)))
}

func TestNotFoundHint(t *testing.T) {
	assert.Equal(t, ui.ThemeNotFoundHint, ui.HintOf(notFoundHint(fmt.Errorf("loading: %w", fs.ErrNotExist))))
	assert.Empty(t, ui.HintOf(notFoundHint(errors.New("other"))))
}
