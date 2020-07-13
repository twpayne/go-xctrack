package xctrack

import (
	"bytes"
	"encoding/json"
)

// ParseTask parses a Task from data.
func ParseTask(data []byte) (*Task, error) {
	if bytes.HasPrefix(data, []byte(QRCodeScheme)) {
		var qrCodeTask QRCodeTask
		if err := json.Unmarshal(data[len(QRCodeScheme):], &qrCodeTask); err != nil {
			return nil, err
		}
		return qrCodeTask.Task(), nil
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}
