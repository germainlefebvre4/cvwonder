package render_html

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"github.com/sirupsen/logrus"
)

// nopWriteCloser wraps a *bytes.Buffer to satisfy io.WriteCloser.
// Close is a no-op since bytes.Buffer needs no cleanup.
type nopWriteCloser struct {
	*bytes.Buffer
}

func (n nopWriteCloser) Close() error { return nil }

// qrOptions holds parsed options for QR code generation with their defaults.
type qrOptions struct {
	size uint8  // pixels per QR cell (default: 5)
	fg   string // foreground hex color (default: #000000)
	bg   string // background hex color (default: #ffffff)
	ec   string // error correction level: L/M/Q/H (default: M)
}

// parseQROptions parses variadic "key=value" option strings into a qrOptions
// struct. Unknown keys are silently ignored. Defaults are applied for any
// key not present in opts.
func parseQROptions(opts []string) qrOptions {
	o := qrOptions{
		size: 5,
		fg:   "#000000",
		bg:   "#ffffff",
		ec:   "M",
	}
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "size":
			if n, err := strconv.ParseUint(val, 10, 8); err == nil {
				o.size = uint8(n)
			}
		case "fg":
			o.fg = val
		case "bg":
			o.bg = val
		case "ec":
			switch strings.ToUpper(val) {
			case "L", "M", "Q", "H":
				o.ec = strings.ToUpper(val)
			}
		}
		// unknown keys are silently ignored
	}
	return o
}

// ecLevel converts the string error correction level to the qrcode constant.
func ecLevel(ec string) qrcode.EncodeOption {
	switch ec {
	case "L":
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionLow)
	case "Q":
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart)
	case "H":
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest)
	default: // "M"
		return qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium)
	}
}

// qrCode generates a QR code from url and returns an HTML <img> tag with the
// PNG image embedded as a base64 data URI. Returns "" if url is empty.
// Supports variadic options: "size=N", "fg=#RRGGBB", "bg=#RRGGBB", "ec=L|M|Q|H".
func qrCode(url string, opts ...string) string {
	if url == "" {
		return ""
	}

	o := parseQROptions(opts)

	qrc, err := qrcode.NewWith(url, ecLevel(o.ec))
	if err != nil {
		logrus.Warnf("qrCode: failed to create QR code for %q: %v", url, err)
		return ""
	}

	buf := &bytes.Buffer{}
	wc := nopWriteCloser{buf}

	w := standard.NewWithWriter(wc,
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithQRWidth(o.size),
		standard.WithFgColorRGBHex(o.fg),
		standard.WithBgColorRGBHex(o.bg),
	)

	if err = qrc.Save(w); err != nil {
		logrus.Warnf("qrCode: failed to save QR code for %q: %v", url, err)
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="QR Code">`, encoded)
}
