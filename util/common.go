package util

import (
	"math"
)

// CalculateDistance returns the distance between points x1,y1 and x2,y2
func CalculateDistance(x1, y1, x2, y2 int) float64 {
	// Use Manhattan distance for grid-based movement
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	return dx + dy
}

// Notifier is an interface that can be implemented to recieve messages
type Notifier interface {
	AddMessage(string)
}
