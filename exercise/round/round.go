package round

import "math"

// Round round a float64
func Round(a float64) float64 {
	return math.Floor(0.5 + a)
}
