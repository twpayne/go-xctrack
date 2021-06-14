package xctrack

import (
	"fmt"
	"regexp"
	"strconv"
)

// Constants.
const (
	TaskExtension = ".xctsk"
	TaskMIMEType  = "application/xctsk"
	TaskVersion   = 1
)

var timeRegexp = regexp.MustCompile(`\A"(\d\d):(\d\d):(\d\d)Z"\z`)

// An Direction is a direction.
type Direction string

// Directions.
const (
	DirectionEnter Direction = "ENTER"
	DirectionExit  Direction = "EXIT"
)

// An EarthModel is an Earth model.
type EarthModel string

// Earth models.
const (
	EarthModelWGS84     EarthModel = "WGS84"
	EarthModelFAISphere EarthModel = "FAI_SPHERE"
)

// A Goal is a goal.
type Goal struct {
	Type     GoalType `json:"type"`
	Deadline *Time    `json:"deadline,omitempty"`
}

// A GoalType is a goal type.
type GoalType string

// Goal types.
const (
	GoalTypeCylinder GoalType = "CYLINDER"
	GoalTypeLine     GoalType = "LINE"
)

// An SSS is a start of speed section.
type SSS struct {
	Type      SSSType   `json:"type"`
	Direction Direction `json:"direction"`
	TimeGates []*Time   `json:"timeGates"`
}

// An SSSType is a start of speed section type.
type SSSType string

// Start of speed section types.
const (
	SSSTypeRace        SSSType = "RACE"
	SSSTypeElapsedTime SSSType = "ELAPSED-TIME"
)

// A Takeoff is a takeoff.
type Takeoff struct {
	TimeOpen  *Time `json:"timeOpen,omitempty"`
	TimeClose *Time `json:"timeClose,omitempty"`
}

// A Task is an XC Track task, see
// http://xctrack.org/Competition_Interfaces.html.
type Task struct {
	TaskType   TaskType     `json:"taskType"`
	Version    int          `json:"version"`
	EarthModel EarthModel   `json:"earthModel,omitempty"`
	Turnpoints []*Turnpoint `json:"turnpoints"`
	Takeoff    *Takeoff     `json:"takeoff,omitempty"`
	SSS        *SSS         `json:"sss,omitempty"`
	Goal       *Goal        `json:"goal,omitempty"`
}

// A TaskType is a task type.
type TaskType string

// Task types.
const (
	TaskTypeClassic TaskType = "CLASSIC"
)

// A Time is a time.
type Time struct {
	Hour   int
	Minute int
	Second int
}

// A Turnpoint is a turnpoint.
type Turnpoint struct {
	Type     TurnpointType `json:"type,omitempty"`
	Radius   int           `json:"radius"`
	Waypoint Waypoint      `json:"waypoint"`
}

// A TurnpointType is a turnpoint type.
type TurnpointType string

// Turnpoint types.
const (
	TurnpointTypeNone    TurnpointType = ""
	TurnpointTypeTakeoff TurnpointType = "TAKEOFF"
	TurnpointTypeSSS     TurnpointType = "SSS"
	TurnpointTypeESS     TurnpointType = "ESS"
)

// A Waypoint is a waypoint.
type Waypoint struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	AltSmoothed int     `json:"altSmoothed"`
}

// An errInvalidTime is an invalid time.
type errInvalidTime string

func (e errInvalidTime) Error() string {
	return fmt.Sprintf("invalid time: %q", string(e))
}

// MarshalJSON implements encoding/json.Marshaler.
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%02d:%02d:%02dZ\"", t.Hour, t.Minute, t.Second)), nil
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (t *Time) UnmarshalJSON(b []byte) error {
	m := timeRegexp.FindSubmatch(b)
	if m == nil {
		return errInvalidTime(b)
	}
	t.Hour, _ = strconv.Atoi(string(m[1]))
	t.Minute, _ = strconv.Atoi(string(m[2]))
	t.Second, _ = strconv.Atoi(string(m[3]))
	return nil
}
