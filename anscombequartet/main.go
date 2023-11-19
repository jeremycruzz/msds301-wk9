package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
)

type dataset struct {
	x []float64
	y []float64
}

func main() {
	anscombeQuartet := []dataset{
		{
			x: []float64{10.0, 8.0, 13.0, 9.0, 11.0, 14.0, 6.0, 4.0, 12.0, 7.0, 5.0},
			y: []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82, 5.68},
		},
		{
			x: []float64{10.0, 8.0, 13.0, 9.0, 11.0, 14.0, 6.0, 4.0, 12.0, 7.0, 5.0},
			y: []float64{9.14, 8.14, 8.74, 8.77, 9.26, 8.10, 6.13, 3.10, 9.13, 7.26, 4.74},
		},
		{
			x: []float64{10.0, 8.0, 13.0, 9.0, 11.0, 14.0, 6.0, 4.0, 12.0, 7.0, 5.0},
			y: []float64{7.46, 6.77, 12.74, 7.11, 7.81, 8.84, 6.08, 5.39, 8.15, 6.42, 5.73},
		},
		{
			x: []float64{8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 8.0, 19.0, 8.0, 8.0, 8.0},
			y: []float64{6.58, 5.76, 7.71, 8.84, 8.47, 7.04, 5.25, 12.50, 5.56, 7.91, 6.89},
		},
	}

	for i, set := range anscombeQuartet {
		alpha, beta := stat.LinearRegression(set.x, set.y, nil, false)
		fmt.Printf("Set %c: m=%.2f b=%.2f\n", 'I'+rune(i), beta, alpha)
	}
}
