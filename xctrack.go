package xctrack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // Parse GIF images.
	_ "image/jpeg" // Parse JPEG images.
	_ "image/png"  // Parse PNG images.
	"io"
	"net/http"

	"github.com/liyue201/goqr"
)

var (
	errEmptyInput    = errors.New("empty input")
	errInvalidFormat = errors.New("invalid format")
)

// LoadTaskFromCode loads a task from a code.
func LoadTaskFromCode(ctx context.Context, code string) (any, error) {
	url := "https://tools.xcontest.org/api/xctsk/load/" + code
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return ParseTask(data)
}

// ParseTask parses a Task from data.
func ParseTask(data []byte) (any, error) {
	switch {
	case len(data) == 0:
		return nil, errEmptyInput
	case bytes.HasPrefix(data, []byte{'{'}):
		var taskType struct {
			TaskType TaskType `json:"taskType"`
		}
		switch err := json.Unmarshal(data, &taskType); {
		case err != nil:
			return nil, err
		case taskType.TaskType == TaskTypeClassic:
			var task Task
			if err := json.Unmarshal(data, &task); err != nil {
				return nil, err
			}
			return &task, nil
		case taskType.TaskType == TaskTypeWaypointList:
			var waypointList WaypointList
			if err := json.Unmarshal(data, &waypointList); err != nil {
				return nil, err
			}
			return &waypointList, nil
		default:
			return nil, fmt.Errorf("%s: unknown task type", taskType.TaskType)
		}
	case bytes.HasPrefix(data, []byte(QRCodeScheme)):
		var qrCodeTask QRCodeTask
		if err := json.Unmarshal(data[len(QRCodeScheme):], &qrCodeTask); err != nil {
			return nil, err
		}
		return qrCodeTask.Task(), nil
	default:
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
}
