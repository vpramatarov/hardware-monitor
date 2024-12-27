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

type HardwareHtml struct{}

func NewHardwareHtml() *HardwareHtml {
	return &HardwareHtml{}
}

func (h *HardwareHtml) GetSystemSection() string {
	runTimeOS := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	hostStat, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}

	html := "<div class='system-data'><table class='table table-striped table-hover table-sm'><tbody>"
	html += "<tr><td>Operating System:</td> <td><i class='fa fa-brands fa-linux'></i> " + runTimeOS + "</td></tr>"
	html += "<tr><td>Platform:</td><td> <i class='fa fa-brands fa-fedora'></i> " + hostStat.Platform + "</td></tr>"
	html += "<tr><td>Hostname:</td><td>" + hostStat.Hostname + "</td></tr>"
	html += "<tr><td>Number of processes running:</td><td>" + strconv.FormatUint(hostStat.Procs, 10) + "</td></tr>"
	html += "<tr><td>Total memory:</td><td>" + FormatUintToUnit(vmStat.Total, "MB") + "</td></tr>"
	html += "<tr><td>Free memory:</td><td>" + FormatUintToUnit(vmStat.Free, "MB") + "</td></tr>"
	html += "<tr><td>Percentage used memory:</td><td>" + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%</td></tr></tbody></table>"
	html += "</div>"

	return html
}

func (h *HardwareHtml) GetDiskSection() string {
	diskStat, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	html := "<div class='disk-data'><table class='table table-striped table-hover table-sm'><tbody>"
	html += "<tr><td>Total disk space:</td><td>" + FormatUintToUnit(diskStat.Total, "GB") + "</td></tr>"
	html += "<tr><td>Used disk space:</td><td>" + FormatUintToUnit(diskStat.Used, "GB") + "</td></tr>"
	html += "<tr><td>Free disk space:</td><td>" + FormatUintToUnit(diskStat.Free, "GB") + "</td></tr>"
	html += "<tr><td>Percentage disk space usage:</td><td>" + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%</td></tr>"
	return html
}

func (h *HardwareHtml) GetCpuSection() string {
	cpuStat, err := cpu.Info()
	if err != nil {
		fmt.Println("Error getting CPU info", err)

	}

	percentage, err := cpu.Percent(0, true)

	if err != nil {
		log.Fatal(err)
	}

	html := "<div class='cpu-data'><table class='table table-striped table-hover table-sm'><tbody>"

	if len(cpuStat) != 0 {
		html += "<tr><td>Model Name:</td><td>" + cpuStat[0].ModelName + "</td></tr>"
		html += "<tr><td>Family:</td><td>" + cpuStat[0].Family + "</td></tr>"
		html += "<tr><td>Speed:</td><td>" + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz</td></tr>"
	}

	firstCpus := percentage[:len(percentage)/2]
	secondCpus := percentage[len(percentage)/2:]
	totalCores := runtime.NumCPU()

	html += "<tr><td>Cores (" + strconv.Itoa(totalCores) + "): </td><td><div class='row mb-4'><div class='col-md-6'><table class='table table-sm'><tbody>"

	for idx, cpupercent := range firstCpus {
		html += "<tr><td>CPU [" + strconv.Itoa(idx+1) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
	}

	html += "</tbody></table></div><div class='col-md-6'><table class='table table-sm'><tbody>"

	secondCpusStartIdx := totalCores / 2

	for idx, cpupercent := range secondCpus {
		html += "<tr><td>CPU [" + strconv.Itoa((idx+1)+secondCpusStartIdx) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
	}

	html += "</tbody></table></div></div></td></tr></tbody></table></div>"

	return html
}
