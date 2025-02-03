// Package movement provides movement strategies for mechs
package movement

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	// Movement constants
	moveStep = 0.1 // Slower movement speed
	directionChangeChance = 0.1
	minStepThreshold = 1.0 // Minimum step size for movement

	// Game boundary constants
	maxLevelWidth = 60
	maxLevelHeight = 40
	minCoordinate = -maxLevelWidth // Allow negative coordinates up to level width
)

// Strategy defines the interface for mech movement behaviors
type Strategy interface {
	// NextMove calculates the next x,y position based on current position
	NextMove(currentX, currentY int) (newX, newY int)
}

// RandomWalkStrategy makes the mech move randomly in any direction
type RandomWalkStrategy struct {
	mu        sync.Mutex
	rng       *rand.Rand
	direction float64
	stepX     float64
	stepY     float64
}

// NewRandomWalkStrategy creates a new random walk movement strategy
func NewRandomWalkStrategy() *RandomWalkStrategy {
	return &RandomWalkStrategy{
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		direction: 0,
		stepX:     0,
		stepY:     0,
	}
}

// updateDirection changes direction with a random chance
func (s *RandomWalkStrategy) updateDirection() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stepX == 0 && s.stepY == 0 || s.rng.Float64() < directionChangeChance {
		s.direction = s.rng.Float64() * 2 * math.Pi
		s.stepX = math.Cos(s.direction) * moveStep
		s.stepY = math.Sin(s.direction) * moveStep
	}
}

// accumulateSteps updates step values based on current direction and returns integer movements
func (s *RandomWalkStrategy) accumulateSteps() (dx, dy int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stepX += math.Cos(s.direction) * moveStep
	s.stepY += math.Sin(s.direction) * moveStep

	// Convert to integer movements
	if math.Abs(s.stepX) >= minStepThreshold {
		dx = int(math.Round(s.stepX))
		s.stepX -= float64(dx)
	}
	if math.Abs(s.stepY) >= minStepThreshold {
		dy = int(math.Round(s.stepY))
		s.stepY -= float64(dy)
	}
	return dx, dy
}

// clampToGameBounds ensures a coordinate stays within game boundaries
func clampToGameBounds(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// NextMove implements Strategy interface
func (s *RandomWalkStrategy) NextMove(currentX, currentY int) (newX, newY int) {
	s.updateDirection()
	dx, dy := s.accumulateSteps()
	
	// Calculate new position
	newX = currentX + dx
	newY = currentY + dy
	
	// Clamp to game boundaries
	newX = clampToGameBounds(newX, minCoordinate, maxLevelWidth)
	newY = clampToGameBounds(newY, minCoordinate, maxLevelHeight)
	
	return newX, newY
}

// PatrolStrategy makes the mech patrol between points
type PatrolStrategy struct {
	points     [][2]int
	currPoint  int
	stepX      float64
	stepY      float64
	targetX    int
	targetY    int
	direction  float64
}

// validatePoint checks if a point is within game boundaries
func validatePoint(x, y int) error {
	if x < minCoordinate || x > maxLevelWidth {
		return fmt.Errorf("x coordinate %d is out of bounds [%d, %d]", x, minCoordinate, maxLevelWidth)
	}
	if y < minCoordinate || y > maxLevelHeight {
		return fmt.Errorf("y coordinate %d is out of bounds [%d, %d]", y, minCoordinate, maxLevelHeight)
	}
	return nil
}

// NewPatrolStrategy creates a new patrol movement strategy
func NewPatrolStrategy(points [][2]int) (*PatrolStrategy, error) {
	if len(points) < 2 {
		return nil, fmt.Errorf("patrol strategy requires at least 2 points, got %d", len(points))
	}

	// Validate all points are within bounds
	for i, point := range points {
		if err := validatePoint(point[0], point[1]); err != nil {
			return nil, fmt.Errorf("patrol point %d is invalid: %w", i, err)
		}
	}

	return &PatrolStrategy{
		points:    points,
		currPoint: 0,
		stepX:     0,
		stepY:     0,
		targetX:   points[0][0],
		targetY:   points[0][1],
	}, nil
}

// updateTarget moves to the next patrol point if current target is reached
func (s *PatrolStrategy) updateTarget(currentX, currentY int) {
	// Check if we've reached the current target
	if currentX == s.targetX && currentY == s.targetY {
		s.currPoint = (s.currPoint + 1) % len(s.points)
		s.targetX = s.points[s.currPoint][0]
		s.targetY = s.points[s.currPoint][1]
	}

	// Calculate direction to target
	dx := float64(s.targetX - currentX)
	dy := float64(s.targetY - currentY)
	s.direction = math.Atan2(dy, dx)
}

// accumulateSteps updates step values based on current direction and returns integer movements
func (s *PatrolStrategy) accumulateSteps() (dx, dy int) {
	s.stepX += math.Cos(s.direction) * moveStep
	s.stepY += math.Sin(s.direction) * moveStep

	// Convert to integer movements
	if math.Abs(s.stepX) >= minStepThreshold {
		dx = int(math.Round(s.stepX))
		s.stepX -= float64(dx)
	}
	if math.Abs(s.stepY) >= minStepThreshold {
		dy = int(math.Round(s.stepY))
		s.stepY -= float64(dy)
	}
	return dx, dy
}

// NextMove implements Strategy interface
func (s *PatrolStrategy) NextMove(currentX, currentY int) (newX, newY int) {
	s.updateTarget(currentX, currentY)
	dx, dy := s.accumulateSteps()
	
	// Calculate new position
	newX = currentX + dx
	newY = currentY + dy
	
	// Clamp to game boundaries
	newX = clampToGameBounds(newX, minCoordinate, maxLevelWidth)
	newY = clampToGameBounds(newY, minCoordinate, maxLevelHeight)
	
	return newX, newY
}
