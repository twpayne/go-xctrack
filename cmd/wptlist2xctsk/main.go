package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/twpayne/go-elevation"
	"github.com/twpayne/go-xctrack"
)

type WptList struct {
	Points []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"points"`
	TaskType string `json:"taskType"`
	Version  int    `json:"version"`
}

func run() error {
	euDEM := flag.String("eu_dem-path", os.Getenv("EU_DEM_PATH"), "path to EU DEM data")
	flag.Parse()

	var wptList WptList
	if err := json.NewDecoder(os.Stdin).Decode(&wptList); err != nil {
		return err
	}

	elevationService, err := elevation.NewEUDEMElevationService(
		os.DirFS(*euDEM),
		elevation.WithCanaryFilename("eu_dem_v11_E40N30.TIF"),
	)
	if err != nil {
		return err
	}

	coords := make([][]float64, len(wptList.Points))
	for i, point := range wptList.Points {
		coords[i] = []float64{point.Lon, point.Lat}
	}
	elevations, err := elevationService.Elevation4326(coords)
	if err != nil {
		return nil
	}
	turnpoints := make([]*xctrack.Turnpoint, 0, len(wptList.Points)+2)
	turnpoints = append(turnpoints, &xctrack.Turnpoint{
		Type:   xctrack.TurnpointTypeTakeoff,
		Radius: 100,
		Waypoint: xctrack.Waypoint{
			Name:        "D01",
			Lat:         wptList.Points[0].Lat,
			Lon:         wptList.Points[0].Lon,
			AltSmoothed: int(elevations[0] + 0.5),
		},
	})
	turnpoints = append(turnpoints, &xctrack.Turnpoint{
		Type:   xctrack.TurnpointTypeSSS,
		Radius: 1000,
		Waypoint: xctrack.Waypoint{
			Name:        "D01",
			Lat:         wptList.Points[0].Lat,
			Lon:         wptList.Points[0].Lon,
			AltSmoothed: int(elevations[0] + 0.5),
		},
	})
	for i := 1; i < len(wptList.Points)-1; i++ {
		turnpoints = append(turnpoints, &xctrack.Turnpoint{
			Radius: 1000,
			Waypoint: xctrack.Waypoint{
				Name:        fmt.Sprintf("B%02d", i),
				Lat:         wptList.Points[i].Lat,
				Lon:         wptList.Points[i].Lon,
				AltSmoothed: int(elevations[i] + 0.5),
			},
		})
	}
	turnpoints = append(turnpoints, &xctrack.Turnpoint{
		Type:   xctrack.TurnpointTypeESS,
		Radius: 1000,
		Waypoint: xctrack.Waypoint{
			Name:        "A01",
			Lat:         wptList.Points[len(wptList.Points)-1].Lat,
			Lon:         wptList.Points[len(wptList.Points)-1].Lon,
			AltSmoothed: int(elevations[len(wptList.Points)-1] + 0.5),
		},
	})
	turnpoints = append(turnpoints, &xctrack.Turnpoint{
		Type:   xctrack.TurnpointTypeESS,
		Radius: 100,
		Waypoint: xctrack.Waypoint{
			Name:        "A01",
			Lat:         wptList.Points[len(wptList.Points)-1].Lat,
			Lon:         wptList.Points[len(wptList.Points)-1].Lon,
			AltSmoothed: int(elevations[len(wptList.Points)-1] + 0.5),
		},
	})
	task := xctrack.Task{
		TaskType:   xctrack.TaskTypeClassic,
		Version:    xctrack.Version,
		Turnpoints: turnpoints,
		Takeoff: &xctrack.Takeoff{
			TimeOpen:  &xctrack.TimeOfDay{Hour: 9},
			TimeClose: &xctrack.TimeOfDay{Hour: 18, Minute: 30},
		},
		SSS: &xctrack.SSS{
			Type:      xctrack.SSSTypeElapsedTime,
			Direction: xctrack.DirectionExit,
			TimeGates: []*xctrack.TimeOfDay{
				{Hour: 9},
			},
		},
		Goal: &xctrack.Goal{
			Type:     xctrack.GoalTypeCylinder,
			Deadline: &xctrack.TimeOfDay{Hour: 18, Minute: 30},
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(task)
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
