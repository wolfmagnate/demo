package main

import "math"

func Ease1(x float64) float64 {
	return -4.14736577*x*x*x*x + 13.54280144*x*x*x - 11.45631207*x*x + 3.0608764*x
}

func Ease2(x float64) float64 {
	if x < 0.33 {
		return math.Sqrt(3*x) / 3
	} else if x < 0.66 {
		x -= 0.33
		return math.Sqrt(3*x)/3 + 0.33
	} else {
		x -= 0.66
		return math.Sqrt(3*x)/3 + 0.66
	}
}

func Ease3(x float64) float64 {
	return x * x * x
}
