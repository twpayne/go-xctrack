package xctrack

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/gif"  // Parse GIF images.
	_ "image/jpeg" // Parse JPEG images.
	_ "image/png"  // Parse PNG images.

	"github.com/liyue201/goqr"
)

var (
	errEmptyInput    = errors.New("empty input")
	errInvalidFormat = errors.New("invalid format")
)

// ParseTask parses a Task from data.
func ParseTask(data []byte) (*Task, error) {
	if len(data) == 0 {
		return nil, errEmptyInput
	}

	if bytes.HasPrefix(data, []byte(QRCodeScheme)) {
		var qrCodeTask QRCodeTask
		if err := json.Unmarshal(data[len(QRCodeScheme):], &qrCodeTask); err == nil {
			return qrCodeTask.Task(), nil
		}
	}

	var task Task
	if err := json.Unmarshal(data, &task); err == nil {
		return &task, nil
	}

	if img, _, err := image.Decode(bytes.NewReader(data)); err == nil {
		if qrCodes, err := goqr.Recognize(img); err == nil {
			for _, qrCode := range qrCodes {
				if bytes.HasPrefix(qrCode.Payload, []byte(QRCodeScheme)) {
					var qrCodeTask QRCodeTask
					if err := json.Unmarshal(qrCode.Payload[len(QRCodeScheme):], &qrCodeTask); err == nil {
						return qrCodeTask.Task(), nil
					}
				}
			}
		}
	}

	return nil, errInvalidFormat
}
