package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/skip2/go-qrcode"

	"github.com/twpayne/go-xctrack"
)

var format = flag.String("format", "json", "format")

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

	switch *format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(task)
	case "qrcode":
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
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
