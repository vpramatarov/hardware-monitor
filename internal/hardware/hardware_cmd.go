package hardware

import (
	"fmt"
	"log"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type HardwareCmd struct{}

func NewHardwareCmd() *HardwareCmd {
	return &HardwareCmd{}
}

func (h *HardwareCmd) GetSystemSection() string {
	runTimeOS := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	hostStat, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}

	output := ""
	output = output + "Operating System: " + runTimeOS + "\n"
	output = output + "Platform: " + hostStat.Platform + "\n"
	output = output + "Hostname: " + hostStat.Hostname + "\n"
	output = output + "Number of processes running: " + strconv.FormatUint(hostStat.Procs, 10) + "\n"
	output = output + "Total memory: " + FormatUintToUnit(vmStat.Total, "MB") + "\n"
	output = output + "Free memory: " + FormatUintToUnit(vmStat.Free, "MB") + "\n"
	output = output + "Percentage used memory: " + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%"

	return output
}

func (h *HardwareCmd) GetDiskSection() string {
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	output := ""
	output = output + "Total disk space: " + FormatUintToUnit(diskStat.Total, "GB") + "\n"
	output = output + "Used disk space: " + FormatUintToUnit(diskStat.Used, "GB") + "\n"
	output = output + "Free disk space: " + FormatUintToUnit(diskStat.Free, "GB") + "\n"
	output = output + "Percentage disk space usage: " + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%"
	return output
}

func (h *HardwareCmd) GetCpuSection() string {
	cpuStat, err := cpu.Info()
	if err != nil {
		fmt.Println("Error getting CPU info", err)
	}

	percentage, err := cpu.Percent(0, true)

	if err != nil {
		log.Fatal(err)
	}

	output := ""

	if len(cpuStat) != 0 {
		output = output + "Model Name: " + cpuStat[0].ModelName + "\n"
		output = output + "Family: " + cpuStat[0].Family + "\n"
		output = output + "Speed: " + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz\n"
	}

	firstCpus := percentage[:len(percentage)/2]
	secondCpus := percentage[len(percentage)/2:]
	totalCores := runtime.NumCPU()

	output = output + "Cores: (" + strconv.Itoa(totalCores) + ")\n"

	for idx, cpupercent := range firstCpus {
		output = output + "CPU [" + strconv.Itoa(idx+1) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%\n"
	}

	secondCpusStartIdx := totalCores / 2

	for idx, cpupercent := range secondCpus {
		output = output + "CPU [" + strconv.Itoa((idx+1)+secondCpusStartIdx) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%\n"
	}

	return output
}
