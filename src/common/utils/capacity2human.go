package utils

import (
	"fmt"
	"math"
	"strconv"
)

// CAP_PRE_UNIT 每单位容量的间隔
const CAP_PRE_UNIT = 1024

// 容量单位表
var capUnitList = []string{"Byte", "KB", "MB", "GB", "TB", "PB", "EB", "BB"}

// Capacity2HumanString 将字节单位的容量转为可读性更高的单位，以字符串形式输出
func Capacity2HumanString(cap uint64, places uint) (finalCapStr string) {
	finalCap, unit := Capacity2Human(cap)

	template := "%." + strconv.Itoa(int(places)) + "f %s"
	return fmt.Sprintf(template, finalCap, unit)
}

// Capacity2Human 将字节单位的容量转为可读性更高的单位
func Capacity2Human(cap uint64) (finalCap float64, unit string) {
	var (
		idx int
		u   string
	)

	defer func() {
		unit = u
	}()

	for idx, u = range capUnitList {
		unit = u
		finalCap = float64(cap) / math.Pow(CAP_PRE_UNIT, float64(idx))
		if finalCap < CAP_PRE_UNIT {
			return
		}
	}

	return
}
