package xctrack

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTask(t *testing.T) {
	for _, tc := range []struct {
		t *Task
		j string
	}{
		{
			t: &Task{
				TaskType: TaskTypeClassic,
				Version:  1,
				Turnpoints: []*Turnpoint{
					&Turnpoint{
						Type:     TurnpointTypeTakeoff,
						Name:     "Takeoff",
						Lat:      1,
						Lon:      2,
						Altitude: 3,
						Radius:   4,
					},
				},
				SSS: &SSS{
					Type:      SSSTypeRace,
					Direction: SSSDirectionEnter,
					TimeGates: []*Time{
						&Time{Hour: 1, Minute: 2, Second: 3},
					},
				},
			},
			j: `{"taskType":"CLASSIC","version":1,"turnpoints":[{"type":"TAKEOFF","name":"Takeoff","lat":1,"lon":2,"altitude":3,"radius":4}],"sss":{"type":"RACE","direction":"ENTER","TimeGates":["01:02:03Z"]}}`,
		},
	} {
		if gotB, gotErr := json.Marshal(tc.t); gotErr != nil || string(gotB) != tc.j {
			t.Errorf("json.Marshal(%#v) == []byte(%q), %v, want []byte(%q), <nil>", tc.t, string(gotB), gotErr, tc.j)
		}
		gotTask := &Task{}
		if gotErr := json.Unmarshal([]byte(tc.j), gotTask); gotErr != nil || !reflect.DeepEqual(gotTask, tc.t) {
			t.Errorf("json.Unmarshal([]byte(%q), ...) == %v and stored %+v, want <nil> and stored %+v", tc.j, gotErr, gotTask, tc.t)
		}
	}
}
