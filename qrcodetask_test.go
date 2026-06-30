package xctrack_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-xctrack/v2"
)

func TestQRCodeTask(t *testing.T) {
	// testdata/peje.task contains a coordinate which, when encoded, includes a backslash
	data, err := os.ReadFile("testdata/peje.xctsk")
	assert.NoError(t, err)
	var task xctrack.Task
	assert.NoError(t, json.Unmarshal(data, &task))
	_, err = task.QRCodeTask().String()
	assert.NoError(t, err)
}

func TestQRCodeTaskZString(t *testing.T) {
	data, err := os.ReadFile("testdata/task_2020-07-07-normalized.xctsk")
	assert.NoError(t, err)
	taskAny, err := xctrack.ParseTask(data)
	assert.NoError(t, err)
	task, ok := taskAny.(*xctrack.Task)
	assert.True(t, ok)
	xctskz, err := task.QRCodeTask().XCTSKZ()
	assert.NoError(t, err)
	parsedTask, err := xctrack.ParseTask([]byte(xctskz))
	assert.NoError(t, err)
	assert.Equal(t, taskAny, parsedTask)
}
