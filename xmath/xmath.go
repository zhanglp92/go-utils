package xmath

/*
#cgo CXXFLAGS: -I. -DLOGGING_LEVEL=LL_WARNING -O3 -Wall
#include <stdlib.h>
#include "xmath.h"
*/
import "C"
import (
	"time"
)

// 自定义基础方法

// 最大/小值函数族 ---------------- begin
// . 约束v的取值范围, mask=00(取值[min, max]), mask=01(取值[min, 无穷)), mask=10(取值(无穷,max])

// BorderLine 取值边界
type BorderLine uint

// 最大/小值边界
const (
	MinBorderLine    BorderLine = 1 << iota // 1<<0
	MaxBorderLine                           // 1<<1
	MinMaxBorderLine = MinBorderLine | MaxBorderLine
)

// MaximinInt ...
func MaximinInt(v *int, min, max int, mask BorderLine) {
	if mask&MaxBorderLine != 0 && *v > max {
		*v = max
	}
	if mask&MinBorderLine != 0 && *v < min {
		*v = min
	}
}

// MaximinInt64 ...
func MaximinInt64(v *int64, min, max int64, mask BorderLine) {
	if mask&MaxBorderLine != 0 && *v > max {
		*v = max
	}
	if mask&MinBorderLine != 0 && *v < min {
		*v = min
	}
}

// MaximinDuration ...
func MaximinDuration(v *time.Duration, min, max time.Duration, mask BorderLine) {
	if mask&MaxBorderLine != 0 && *v > max {
		*v = max
	}
	if mask&MinBorderLine != 0 && *v < min {
		*v = min
	}
}

// 最大/小值函数族 ---------------- end

// callC c调用测试
func callC() {
	C.call_c_test()
}
