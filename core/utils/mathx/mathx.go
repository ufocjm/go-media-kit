package mathx

import (
	"math"
)

func Round(x float64, precision int) float64 {
	scale := math.Pow(10, float64(precision))
	return math.Round(x*scale) / scale
}
