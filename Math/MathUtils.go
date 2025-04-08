package Math

import "math"

func DegToRad(deg float64) float64 {
	return deg / 180 * math.Pi
}

func RadToDeg(rad float64) float64 {
	return rad / math.Pi * 180
}
