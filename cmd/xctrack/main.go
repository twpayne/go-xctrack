package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/skip2/go-qrcode"
	"github.com/twpayne/go-kml/v3"

	"github.com/twpayne/go-xctrack"
)

var format = flag.String("format", "json", "format")

func taskToKML(t *xctrack.Task) *kml.KMLElement {
	coordinates := make([]kml.Coordinate, len(t.Turnpoints))
	for i, turnpoint := range t.Turnpoints {
		coordinates[i] = kml.Coordinate{
			Lat: turnpoint.Waypoint.Lat,
			Lon: turnpoint.Waypoint.Lon,
			Alt: float64(turnpoint.Waypoint.AltSmoothed),
		}
	}
	return kml.KML(
		kml.Folder(
			kml.Placemark(
				kml.LineString(
					kml.Coordinates(coordinates...),
				),
			),
		),
	)
}

func waypointListToKML(l *xctrack.WaypointList) *kml.KMLElement {
	coordinates := make([]kml.Coordinate, len(l.Points))
	for i, waypoint := range l.Points {
		coordinates[i] = kml.Coordinate{
			Lat: waypoint.Lat,
			Lon: waypoint.Lon,
			Alt: float64(waypoint.AltSmoothed),
		}
	}
	return kml.KML(
		kml.Folder(
			kml.Placemark(
				kml.LineString(
					kml.Coordinates(coordinates...),
				),
			),
		),
	)
}

func run() error {
	flag.Parse()

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	task, err := xctrack.ParseTask(data)
	if err != nil {
		return err
	}

	switch task := task.(type) {
	case *xctrack.Task:
		switch *format {
		case "kml":
			return taskToKML(task).WriteIndent(os.Stdout, "", "  ")
		case "json":
			return json.NewEncoder(os.Stdout).Encode(task)
		case "png":
			s, err := task.QRCodeTask().String()
			if err != nil {
				return err
			}
			png, err := qrcode.Encode(s, qrcode.Medium, 1024)
			if err != nil {
				return err
			}
			_, err = os.Stdout.Write(png)
			return err
		case "qrcode-json":
			s, err := task.QRCodeTask().String()
			if err != nil {
				return err
			}
			_, err = os.Stdout.WriteString(s)
			return err
		default:
			return fmt.Errorf("%s: invalid format", *format)
		}
	case *xctrack.WaypointList:
		switch *format {
		case "kml":
			return waypointListToKML(task).WriteIndent(os.Stdout, "", "  ")
		case "json":
			return json.NewEncoder(os.Stdout).Encode(task)
		default:
			return fmt.Errorf("%s: invalid format", *format)
		}
	default:
		return fmt.Errorf("%T: unsupported type", task)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
