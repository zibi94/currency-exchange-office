package num

import "math"

func Round(num float64, n int) float64 {
	m := math.Pow10(n)
	return math.Round(num*m) / m
}
