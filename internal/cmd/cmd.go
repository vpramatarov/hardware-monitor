package cmd

import (
	"fmt"
	"time"

	"github.com/vpramatarov/hardware-monitor/internal/hardware"
)

type Cmd struct{}

func NewCmd() *Cmd {
	return &Cmd{}
}

func (h *Cmd) DisplaySystemData() {
	fmt.Println("Start hardware monitoring...")

	hardwareCmd := hardware.NewHardwareCmd()

	go func(h *hardware.HardwareCmd) {
		for {
			timeStamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Println(timeStamp)
			fmt.Println("=================")
			fmt.Println(h.GetSystemSection())
			fmt.Println("=================")
			fmt.Println(h.GetDiskSection())
			fmt.Println("=================")
			fmt.Println(h.GetCpuSection())
			fmt.Println("=================")
			time.Sleep(time.Duration(hardware.SecondsInterval) * time.Second)
		}
	}(hardwareCmd)

	// sleep the main thread
	time.Sleep(time.Duration(hardware.SecondsInterval) * time.Minute)
}
