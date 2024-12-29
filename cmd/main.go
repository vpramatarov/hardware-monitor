package main

import (
	"os"

	"github.com/vpramatarov/hardware-monitor/internal/cmd"
	"github.com/vpramatarov/hardware-monitor/internal/hardware"
	"github.com/vpramatarov/hardware-monitor/internal/ws"
)

func main() {
	display := "cmd" // default

	if len(os.Args) > 1 {
		if os.Args[1] == "cmd" || os.Args[1] == "ws" {
			display = os.Args[1]
		}
	}

	start(display)
}

func start(display string) {
	var runner hardware.IDisplay

	if display == "ws" {
		runner = ws.NewWs()
	} else {
		runner = cmd.NewCmd()
	}

	runner.DisplaySystemData()
}
