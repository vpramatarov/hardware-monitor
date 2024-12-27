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

	run(display)
}

func run(display string) {
	var runner hardware.IDisplay

	if display == "cmd" {
		runner = cmd.NewCmd()
	} else if display == "ws" {
		runner = ws.NewWs()
	} else {
		panic("Invalid runner. Must be `cmd` or `ws`. Default is `ws`.")
	}

	show(runner)
}

func show(hardware hardware.IDisplay) {
	hardware.DisplaySystemData()
}
