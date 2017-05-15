package xctrack

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	TaskExtension    = ".xctsk"
	TaskMIMEType     = "application/xctsk"
	TaskQRCodeScheme = "XCTSK"
)

var (
	timeRegexp = regexp.MustCompile(`\A"(\d\d):(\d\d):(\d\d)Z"\z`)
)

type ESS struct {
	AltitudeTimeBonus float64 `json:"altitudeTimeBonus,omitempty"`
}

type Goal struct {
	Type     GoalType `json:"goalType"`
	Deadline *Time    `json:"deadline"`
}

type GoalType string

const (
	GoalTypeCylinder GoalType = "CYLINDER"
	GoalTypeLine     GoalType = "LINE"
)

type SSS struct {
	Type          SSSType      `json:"type"`
	Direction     SSSDirection `json:"direction"`
	TimeGates     []*Time      `json:"timeGates"`
	TimeClose     *Time        `json:"timeClose,omitempty"`
	TimeLastStart *Time        `json:"timeLastStart,omitempty"`
}

type SSSDirection string

const (
	SSSDirectionEnter SSSDirection = "ENTER"
	SSSDirectionExit  SSSDirection = "EXIT"
)

type SSSType string

const (
	SSSTypeRace        SSSType = "RACE"
	SSSTypeElapsedType SSSType = "ELAPSED-TIME"
)

type Takeoff struct {
	TimeOpen  *Time `json:"timeOpen,omitempty"`
	TimeClose *Time `json:"timeClose,omitempty"`
}

// A Task is an XC Track task, see
// http://xctrack.org/Competition_Interfaces.html.
type Task struct {
	TaskType   TaskType     `json:"taskType"`
	Version    int          `json:"version"`
	Turnpoints []*Turnpoint `json:"turnpoints"`
	Takeoff    *Takeoff     `json:"takeoff,omitempty"`
	SSS        *SSS         `json:"sss"`
	ESS        *ESS         `json:"ess,omitempty"`
}

type TaskType string

const (
	TaskTypeClassic TaskType = "CLASSIC"
)

type Time struct {
	Hour   int
	Minute int
	Second int
}

type Turnpoint struct {
	Type        TurnpointType `json:"type"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Lat         float64       `json:"lat"`
	Lon         float64       `json:"lon"`
	Altitude    float64       `json:"altitude"`
	Radius      float64       `json:"radius"`
}

type TurnpointType string

const (
	TurnpointTypeTakeoff TurnpointType = "TAKEOFF"
	TurnpointTypeSSS     TurnpointType = "SSS"
	TurnpointTypeESS     TurnpointType = "ESS"
)

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%02d:%02d:%02dZ\"", t.Hour, t.Minute, t.Second)), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	m := timeRegexp.FindSubmatch(b)
	if m == nil {
		return fmt.Errorf("invalid time format: %q", string(b))
	}
	var err error
	t.Hour, err = strconv.Atoi(string(m[1]))
	if err != nil {
		panic(err)
	}
	t.Minute, err = strconv.Atoi(string(m[2]))
	if err != nil {
		panic(err)
	}
	t.Second, err = strconv.Atoi(string(m[3]))
	if err != nil {
		panic(err)
	}
	return nil
}
