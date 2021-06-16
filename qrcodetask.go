package xctrack

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	polyline "github.com/twpayne/go-polyline"
)

// Constants.
const (
	QRCodeScheme      = "XCTSK:"
	QRCodeTaskVersion = 2
)

// A QRCodeDirection is a QR code direction.
type QRCodeDirection int

// QR code directions.
const (
	QRCodeDirectionEnter QRCodeDirection = 1
	QRCodeDirectionExit  QRCodeDirection = 2
)

// A QRCodeEarthModel is a QR code Earth model.
type QRCodeEarthModel int

// QR code Earth models.
const (
	QRCodeEarthModelWGS84     QRCodeEarthModel = 0
	QRCodeEarthModelFAISphere QRCodeEarthModel = 1
)

// A QRCodeGoal is a QR code goal.
type QRCodeGoal struct {
	Deadline *Time          `json:"d,omitempty"`
	Type     QRCodeGoalType `json:"t,omitempty"`
}

// A QRCodeGoalType is a QR code goal type.
type QRCodeGoalType int

// QR code goal types.
const (
	QRCodeGoalTypeLine     QRCodeGoalType = 1
	QRCodeGoalTypeCylinder QRCodeGoalType = 2
)

// A QRCodeSSSType is a QR code start of speed section type.
type QRCodeSSSType int

// QR code start of speed section types.
const (
	QRCodeSSSTypeRace        QRCodeSSSType = 1
	QRCodeSSSTypeElapsedTime QRCodeSSSType = 2
)

// A QRCodeSSS is a QR code start.
type QRCodeSSS struct {
	TimeGates []*Time         `json:"g"`
	Direction QRCodeDirection `json:"d"`
	Type      QRCodeSSSType   `json:"t"`
}

// A QRCodeTask is a QR code task.
type QRCodeTask struct {
	TaskType   TaskType           `json:"taskType"`
	Version    int                `json:"version"`
	Turnpoints []*QRCodeTurnpoint `json:"t"`
	SSS        *QRCodeSSS         `json:"s,omitempty"`
	Goal       *QRCodeGoal        `json:"g,omitempty"`
	EarthModel QRCodeEarthModel   `json:"e"`
}

// A QRCodeTaskType is a QR code task type.
type QRCodeTaskType int

// QR code task types.
const (
	QRCodeTaskTypeRace        QRCodeTaskType = 1
	QRCodeTaskTypeElapsedTime QRCodeTaskType = 2
)

// A QRCodeTurnpoint is a QR code turnpoint.
type QRCodeTurnpoint struct {
	Z           QRCodeTurnpointZ    `json:"z"`
	Name        string              `json:"n"`
	Description string              `json:"d,omitempty"`
	Type        QRCodeTurnpointType `json:"t,omitempty"`
}

// A QRCodeTurnpointType is a QR code turnpoint type.
type QRCodeTurnpointType int

// QR code turnpoint types.
const (
	QRCodeTurnpointTypeNone QRCodeTurnpointType = 0
	QRCodeTurnpointTypeSSS  QRCodeTurnpointType = 2
	QRCodeTurnpointTypeESS  QRCodeTurnpointType = 3
)

// A QRCodeTurnpointZ is a QR code turnpoint Z.
type QRCodeTurnpointZ struct {
	Lon    float64
	Lat    float64
	Alt    int
	Radius int
}

var (
	errExpectedClosingDoubleQuote = errors.New("expected closing double quote")
	errExpectedOpeningDoubleQuote = errors.New("expected opening double quote")
	errTrailingBytes              = errors.New("trailing bytes")
)

// An errInvalidQRCodeTurnpointZ is an invalid QR code turnpoint Z error.
type errInvalidQRCodeTurnpointZ struct {
	Value []byte
	Err   error
}

func (e errInvalidQRCodeTurnpointZ) Error() string {
	return fmt.Sprintf("invalid QR code turnpoint z: %q: %v", string(e.Value), e.Err)
}

var (
	qrCodeDirectionValue = map[Direction]QRCodeDirection{
		DirectionEnter: QRCodeDirectionEnter,
		DirectionExit:  QRCodeDirectionExit,
	}
	qrCodeEarthModelValue = map[EarthModel]QRCodeEarthModel{
		EarthModelWGS84:     QRCodeEarthModelWGS84,
		EarthModelFAISphere: QRCodeEarthModelFAISphere,
	}
	qrCodeGoalTypeValue = map[GoalType]QRCodeGoalType{
		GoalTypeLine:     QRCodeGoalTypeLine,
		GoalTypeCylinder: QRCodeGoalTypeCylinder,
	}
	qrCodeSSSTypeValue = map[SSSType]QRCodeSSSType{
		SSSTypeRace:        QRCodeSSSTypeRace,
		SSSTypeElapsedTime: QRCodeSSSTypeElapsedTime,
	}
	qrCodeTurnpointTypeValue = map[TurnpointType]QRCodeTurnpointType{
		TurnpointTypeNone: QRCodeTurnpointTypeNone,
		TurnpointTypeSSS:  QRCodeTurnpointTypeSSS,
		TurnpointTypeESS:  QRCodeTurnpointTypeESS,
	}

	directionValue = map[QRCodeDirection]Direction{
		QRCodeDirectionEnter: DirectionEnter,
		QRCodeDirectionExit:  DirectionExit,
	}
	earthModelValue = map[QRCodeEarthModel]EarthModel{
		QRCodeEarthModelWGS84:     EarthModelWGS84,
		QRCodeEarthModelFAISphere: EarthModelFAISphere,
	}
	goalTypeValue = map[QRCodeGoalType]GoalType{
		QRCodeGoalTypeLine:     GoalTypeLine,
		QRCodeGoalTypeCylinder: GoalTypeCylinder,
	}
	sssTypeValue = map[QRCodeSSSType]SSSType{
		QRCodeSSSTypeRace:        SSSTypeRace,
		QRCodeSSSTypeElapsedTime: SSSTypeElapsedTime,
	}
	turnpointTypeValue = map[QRCodeTurnpointType]TurnpointType{
		QRCodeTurnpointTypeNone: TurnpointTypeNone,
		QRCodeTurnpointTypeSSS:  TurnpointTypeSSS,
		QRCodeTurnpointTypeESS:  TurnpointTypeESS,
	}
)

// MarshalJSON implements encoding/json.Marshaler.
func (z *QRCodeTurnpointZ) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 64)
	buf = append(buf, '"')
	buf = polyline.EncodeInt(buf, int(math.Round(1e5*z.Lon)))
	buf = polyline.EncodeInt(buf, int(math.Round(1e5*z.Lat)))
	buf = polyline.EncodeInt(buf, z.Alt)
	buf = polyline.EncodeInt(buf, z.Radius)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON implements encoding/json.Unmarshaler.
func (z *QRCodeTurnpointZ) UnmarshalJSON(value []byte) error {
	b := value
	if len(b) == 0 || b[0] != '"' {
		return &errInvalidQRCodeTurnpointZ{
			Value: value,
			Err:   errExpectedOpeningDoubleQuote,
		}
	}
	b = b[1:]
	lon, b, err := polyline.DecodeInt(b)
	if err != nil {
		return &errInvalidQRCodeTurnpointZ{
			Value: value,
			Err:   err,
		}
	}
	lat, b, err := polyline.DecodeInt(b)
	if err != nil {
		return &errInvalidQRCodeTurnpointZ{
			Value: value,
			Err:   err,
		}
	}
	alt, b, err := polyline.DecodeInt(b)
	if err != nil {
		return &errInvalidQRCodeTurnpointZ{
			Value: value,
			Err:   err,
		}
	}
	radius, b, err := polyline.DecodeInt(b)
	if err != nil {
		return &errInvalidQRCodeTurnpointZ{
			Value: value,
			Err:   err,
		}
	}
	if len(b) == 0 || b[0] != '"' {
		return errInvalidQRCodeTurnpointZ{
			Value: value,
			Err:   errExpectedClosingDoubleQuote,
		}
	}
	b = b[1:]
	if len(b) != 0 {
		return &errInvalidQRCodeTurnpointZ{
			Value: value,
			Err:   errTrailingBytes,
		}
	}
	z.Lon = float64(lon) / 1e5
	z.Lat = float64(lat) / 1e5
	z.Alt = alt
	z.Radius = radius
	return nil
}

// QRCodeTask returns t as a QRCodeTask.
func (t *Task) QRCodeTask() *QRCodeTask {
	qrCodeTask := &QRCodeTask{
		TaskType:   t.TaskType,
		Version:    QRCodeTaskVersion,
		Turnpoints: make([]*QRCodeTurnpoint, 0, len(t.Turnpoints)),
		EarthModel: qrCodeEarthModelValue[t.EarthModel],
	}
	for _, turnpoint := range t.Turnpoints {
		qrCodeTurnpoint := &QRCodeTurnpoint{
			Z: QRCodeTurnpointZ{
				Lon:    turnpoint.Waypoint.Lon,
				Lat:    turnpoint.Waypoint.Lat,
				Alt:    turnpoint.Waypoint.AltSmoothed,
				Radius: turnpoint.Radius,
			},
			Name:        turnpoint.Waypoint.Name,
			Description: turnpoint.Waypoint.Description,
			Type:        qrCodeTurnpointTypeValue[turnpoint.Type],
		}
		qrCodeTask.Turnpoints = append(qrCodeTask.Turnpoints, qrCodeTurnpoint)
	}
	if t.SSS != nil {
		qrCodeTask.SSS = &QRCodeSSS{
			TimeGates: t.SSS.TimeGates,
			Direction: qrCodeDirectionValue[t.SSS.Direction],
			Type:      qrCodeSSSTypeValue[t.SSS.Type],
		}
	}
	if t.Goal != nil {
		qrCodeTask.Goal = &QRCodeGoal{
			Deadline: t.Goal.Deadline,
			Type:     qrCodeGoalTypeValue[t.Goal.Type],
		}
	}
	return qrCodeTask
}

func (q *QRCodeTask) String() (string, error) {
	sb := &strings.Builder{}
	sb.Grow(4096)
	sb.WriteString(QRCodeScheme)
	if err := json.NewEncoder(sb).Encode(q); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// Task returns t as a Task.
func (q *QRCodeTask) Task() *Task {
	task := &Task{
		TaskType:   q.TaskType,
		Version:    TaskVersion,
		Turnpoints: make([]*Turnpoint, 0, len(q.Turnpoints)),
		EarthModel: earthModelValue[q.EarthModel],
	}
	for _, qrCodeTurnpoint := range q.Turnpoints {
		turnpoint := &Turnpoint{
			Type:   turnpointTypeValue[qrCodeTurnpoint.Type],
			Radius: qrCodeTurnpoint.Z.Radius,
			Waypoint: Waypoint{
				Name:        qrCodeTurnpoint.Name,
				Description: qrCodeTurnpoint.Description,
				Lat:         qrCodeTurnpoint.Z.Lat,
				Lon:         qrCodeTurnpoint.Z.Lon,
				AltSmoothed: qrCodeTurnpoint.Z.Alt,
			},
		}
		task.Turnpoints = append(task.Turnpoints, turnpoint)
	}
	if q.SSS != nil {
		task.SSS = &SSS{
			Type:      sssTypeValue[q.SSS.Type],
			Direction: directionValue[q.SSS.Direction],
			TimeGates: q.SSS.TimeGates,
		}
	}
	if q.Goal != nil {
		task.Goal = &Goal{
			Type:     goalTypeValue[q.Goal.Type],
			Deadline: q.Goal.Deadline,
		}
	}
	return task
}
