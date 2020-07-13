package xctrack_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xctrack "github.com/twpayne/go-xctrack"
)

func TestTask(t *testing.T) {
	for i, tc := range []struct {
		task      *xctrack.Task
		jsonStr   string
		qrCodeStr string
	}{
		{
			task: &xctrack.Task{
				TaskType:   xctrack.TaskTypeClassic,
				Version:    xctrack.TaskVersion,
				EarthModel: xctrack.EarthModelWGS84,
				Turnpoints: []*xctrack.Turnpoint{
					{
						Radius: 1,
						Waypoint: xctrack.Waypoint{
							Name:        "D01",
							Lat:         1,
							Lon:         2,
							AltSmoothed: 3,
						},
					},
				},
				SSS: &xctrack.SSS{
					Type:      xctrack.SSSTypeRace,
					Direction: xctrack.DirectionEnter,
					TimeGates: []*xctrack.Time{
						{Hour: 1, Minute: 2, Second: 3},
					},
				},
				Goal: &xctrack.Goal{
					Type: xctrack.GoalTypeLine,
				},
			},
			jsonStr:   `{"taskType":"CLASSIC","version":1,"earthModel":"WGS84","turnpoints":[{"radius":1,"waypoint":{"name":"D01","lat":1,"lon":2,"altSmoothed":3}}],"sss":{"type":"RACE","direction":"ENTER","timeGates":["01:02:03Z"]},"goal":{"type":"LINE"}}`,
			qrCodeStr: `XCTSK:{"taskType":"CLASSIC","version":2,"t":[{"z":"_seK_ibEEA","n":"D01"}],"s":{"g":["01:02:03Z"],"d":1,"t":1},"g":{"t":1},"e":0}`,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := json.Marshal(tc.task)
			require.NoError(t, err)
			assert.Equal(t, tc.jsonStr, string(b))

			actualQRCodeStr, err := tc.task.QRCodeTask().String()
			require.NoError(t, err)
			assert.Equal(t, tc.qrCodeStr, strings.TrimSpace(actualQRCodeStr))

			actualTask, err := xctrack.ParseTask([]byte(tc.jsonStr))
			require.NoError(t, err)
			assert.Equal(t, tc.task, actualTask)

			actualTask, err = xctrack.ParseTask([]byte(tc.qrCodeStr))
			require.NoError(t, err)
			assert.Equal(t, tc.task, actualTask)
		})
	}
}

func TestTestData(t *testing.T) {
	infos, err := ioutil.ReadDir("testdata")
	require.NoError(t, err)
	for _, info := range infos {
		if !strings.HasSuffix(info.Name(), xctrack.TaskExtension) {
			continue
		}
		t.Run(info.Name(), func(t *testing.T) {
			data, err := ioutil.ReadFile(filepath.Join("testdata", info.Name()))
			require.NoError(t, err)
			_, err = xctrack.ParseTask(data)
			require.NoError(t, err)
		})
	}
}
