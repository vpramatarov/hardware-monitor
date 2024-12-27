package hardware

import (
	"strconv"
)

const megabyteDiv uint64 = 1024 * 1024
const gigabyteDiv uint64 = megabyteDiv * 1024
const SecondsInterval = 5

type IHardware interface {
	GetSystemSection() string
	GetDiskSection() string
	GetCpuSection() string
}

type IDisplay interface {
	DisplaySystemData()
}

func FormatUintToUnit(i uint64, unit string) string {
	if unit == "MB" {
		return strconv.FormatUint(i/megabyteDiv, 10) + " MB"
	} else if unit == "GB" {
		return strconv.FormatUint(i/gigabyteDiv, 10) + " GB"
	} else {
		panic("Unknown unit. Must be MB or GB")
	}
}
