package ui

import (
	"fmt"
	"time"
)

// FormatDuration renders a duration compactly: "<1ms", "148ms", "1.9s", "1m05s".
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return "<1ms"
	}
	if ms := d.Round(time.Millisecond); ms < time.Second {
		return fmt.Sprintf("%dms", ms.Milliseconds())
	}
	if tenths := d.Round(100 * time.Millisecond); tenths < time.Minute {
		return fmt.Sprintf("%.1fs", tenths.Seconds())
	}
	secs := int(d.Round(time.Second).Seconds())
	return fmt.Sprintf("%dm%02ds", secs/60, secs%60)
}
