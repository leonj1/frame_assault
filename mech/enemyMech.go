package mech

import (
	"github.com/Ariemeth/frame_assault/mech/movement"
	tl "github.com/Ariemeth/termloop"
)

const (
	// moveDelayTicks represents how many ticks to wait between moves
	// Since we're running at 2 FPS, setting this to 4 means moving every 2 seconds
	moveDelayTicks = 4
)

// EnemyMech represents an autonomous enemy mech
type EnemyMech struct {
	*Mech
	moveStrategy movement.Strategy
	moveDelay   int
	tickCount   int
}

// NewEnemyMech creates a new enemy mech instance
func NewEnemyMech(name string, maxStructure, x, y int, color tl.Attr, symbol rune, strategy movement.Strategy) *EnemyMech {
	return &EnemyMech{
		Mech:         NewMech(name, maxStructure, x, y, color, symbol),
		moveStrategy: strategy,
		moveDelay:    moveDelayTicks,
		tickCount:    0,
	}
}

// Tick handles the enemy mech's autonomous behavior
func (e *EnemyMech) Tick(event tl.Event) {
	// Call base Mech's Tick first
	e.Mech.Tick(event)

	// Only move if the mech is not destroyed
	if !e.IsDestroyed() {
		if e.game != nil {
			e.game.Log("Enemy %s tick: count=%d", e.Name(), e.tickCount)
		}
		
		e.tickCount++
		if e.tickCount >= e.moveDelay {
			e.tickCount = 0
			
			// Get current position
			currentX, currentY := e.Position()
			
			// Get next move from strategy
			newX, newY := e.moveStrategy.NextMove(currentX, currentY)

			// Validate move before applying
			if !e.isValidMove(newX, newY) {
				return
			}

			if e.game != nil {
				e.game.Log("Enemy %s moving from (%d,%d) to (%d,%d)", 
					e.Name(), currentX, currentY, newX, newY)
			}
			
			// Store current position as previous
			e.prevX, e.prevY = currentX, currentY
			
			// Update position
			e.entity.SetPosition(newX, newY)
		}
	}
}
