package tools

import (
	"fmt"
	"math"
	"strconv"
)

func DecimalForString(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func DecimalForMath(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}
