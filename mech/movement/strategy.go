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

// NextMove implements Strategy interface
func (s *RandomWalkStrategy) NextMove(currentX, currentY int) (newX, newY int) {
	s.updateDirection()
	dx, dy := s.accumulateSteps()
	return currentX + dx, currentY + dy
}

// PatrolStrategy makes the mech patrol between points
type PatrolStrategy struct {
	points     [][2]int
	currPoint  int
	stepX      float64
	stepY      float64
	targetX    int
	targetY    int
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
func NewPatrolStrategy(points [][2]int) *PatrolStrategy {
	if len(points) < 2 {
		panic("PatrolStrategy requires at least 2 points")
	}

	// Validate all points are within bounds
	for i, point := range points {
		if err := validatePoint(point[0], point[1]); err != nil {
			panic(fmt.Sprintf("patrol point %d is invalid: %v", i, err))
		}
	}

	return &PatrolStrategy{
		points:    points,
		currPoint: 0,
		stepX:     0,
		stepY:     0,
		targetX:   points[0][0],
		targetY:   points[0][1],
	}
}

// updateTarget moves to the next patrol point if current target is reached
func (s *PatrolStrategy) updateTarget(currentX, currentY int) {
	if currentX == s.targetX && currentY == s.targetY {
		s.currPoint = (s.currPoint + 1) % len(s.points)
		s.targetX = s.points[s.currPoint][0]
		s.targetY = s.points[s.currPoint][1]
	}
}

// calculateSteps accumulates movement steps toward target
func (s *PatrolStrategy) calculateSteps(currentX, currentY int) (moveX, moveY int) {
	dx := float64(s.targetX - currentX)
	dy := float64(s.targetY - currentY)

	dx, dy = s.normalizeDirection(dx, dy)
	s.accumulateSteps(dx, dy)
	return s.getIntegerSteps()
}

// normalizeDirection normalizes the direction vector without using square root
func (s *PatrolStrategy) normalizeDirection(dx, dy float64) (float64, float64) {
	dxAbs, dyAbs := math.Abs(dx), math.Abs(dy)
	if dxAbs == 0 && dyAbs == 0 {
		return 0, 0
	}

	if dxAbs > dyAbs {
		return dx/dxAbs, dy/dxAbs
	}
	return dx/dyAbs, dy/dyAbs
}

// accumulateSteps adds the normalized direction to the current steps
func (s *PatrolStrategy) accumulateSteps(dx, dy float64) {
	s.stepX += dx * moveStep
	s.stepY += dy * moveStep
}

// getIntegerSteps converts accumulated steps to integer movements
func (s *PatrolStrategy) getIntegerSteps() (moveX, moveY int) {
	if math.Abs(s.stepX) >= minStepThreshold {
		moveX = int(math.Round(s.stepX))
		s.stepX -= float64(moveX)
	}
	if math.Abs(s.stepY) >= minStepThreshold {
		moveY = int(math.Round(s.stepY))
		s.stepY -= float64(moveY)
	}
	return moveX, moveY
}

// NextMove implements Strategy interface
func (s *PatrolStrategy) NextMove(currentX, currentY int) (newX, newY int) {
	s.updateTarget(currentX, currentY)
	moveX, moveY := s.calculateSteps(currentX, currentY)
	return currentX + moveX, currentY + moveY
}
