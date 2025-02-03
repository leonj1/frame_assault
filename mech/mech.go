// Package mech project mech.go
package mech

import (
	"strconv"

	"github.com/Ariemeth/frame_assault/mech/weapon"
	"github.com/Ariemeth/frame_assault/util"
	tl "github.com/Ariemeth/termloop"
)

// Mech is a basic mech type
type Mech struct {
	structure    int
	maxStructure int
	weapons      []weapon.Weapon
	name         string
	entity       *tl.Entity
	prevX        int
	prevY        int
	game         *tl.Game
	level        *tl.BaseLevel
	notifier     util.Notifier
}

const (
	// Game boundary constants
	maxLevelWidth = 60
	maxLevelHeight = 40
	minCoordinate = -maxLevelWidth // Allow negative coordinates up to level width
)

// NewMech is used to create a new instance of a mech with default structure.
func NewMech(name string, maxStructure, x, y int, color tl.Attr, symbol rune) *Mech {
	newMech := Mech{
		name:         name,
		structure:    maxStructure,
		maxStructure: maxStructure,
		entity:       tl.NewEntity(x, y, 1, 1),
	}

	newMech.entity.SetCell(0, 0, &tl.Cell{Fg: color, Ch: symbol})
	return &newMech
}

// AttachGame is used to attach the termloop game struct for logging
func (m *Mech) AttachGame(game *tl.Game) {
	m.game = game
}

// SetLevel sets the game level for the mech
func (m *Mech) SetLevel(level *tl.BaseLevel) {
	m.level = level
	// Update all weapons with the new level
	for i := range m.weapons {
		m.weapons[i].SetLevel(level)
	}
}

// AttachNotifier is used to attach a notification display
func (m *Mech) AttachNotifier(notifier util.Notifier) {
	m.notifier = notifier
}

// Name returns the name of the mech
func (m Mech) Name() string {
	return m.name
}

// Weapons returns the mechs weapons
func (m Mech) Weapons() []weapon.Weapon {
	return m.weapons
}

// StructureLeft Retrieves the amount of remaining structure a mech has.
func (m Mech) StructureLeft() int {
	return m.structure
}

// Size returns the height and width of the mech
func (m Mech) Size() (int, int) {
	return m.entity.Size()
}

// Position returns the x,y position of the mech
func (m Mech) Position() (int, int) {
	return m.entity.Position()
}

// Collide is used called to see if the mech collided with another physical object
func (m *Mech) Collide(collision tl.Physical) {
	// Check if it's a Rectangle we're colliding with
	if _, ok := collision.(*tl.Rectangle); ok {
		m.entity.SetPosition(m.prevX, m.prevY)
		// or if it is another mech
	} else if _, ok := collision.(*Mech); ok {
		m.entity.SetPosition(m.prevX, m.prevY)
	}
}

// Draw passes the draw call to entity.
func (m *Mech) Draw(screen *tl.Screen) {
	if m.StructureLeft() > 0 {
		m.entity.Draw(screen)
	}
}

// Tick is called to process 1 tick of actions based on the
// type of event.
func (m *Mech) Tick(event tl.Event) {
	m.prevX, m.prevY = m.entity.Position()

	// Update level reference if needed
	if m.level == nil && m.game != nil && m.game.Screen() != nil {
		if level, ok := m.game.Screen().Level().(*tl.BaseLevel); ok {
			m.SetLevel(level)
		}
	}
}

// Hit is call when a mech is hit
func (m *Mech) Hit(damage int) {
	//check if the mech is already destroyed
	if m.structure <= 0 {
		return
	}

	m.structure -= damage
	message1 := m.name + " takes " + strconv.Itoa(damage)
	m.game.Log(message1)
	m.notifier.AddMessage(message1)

	if m.structure <= 0 {
		message2 := m.name + " has been destroyed"
		m.game.Log(message2)
		m.notifier.AddMessage(message2)
		m.game.Screen().Level().RemoveEntity(m)
	}
}

// IsDestroyed returns true is the target is destroyed, false otherwise.
func (m Mech) IsDestroyed() bool {
	return m.structure <= 0
}

// AddWeapon adds a Weapon to the mech
func (m *Mech) AddWeapon(w weapon.Weapon) {
	// Set the weapon's level for bullet creation if we have one
	if m.level != nil {
		w.SetLevel(m.level)
	}
	m.weapons = append(m.weapons, w)
}

// Fire tells the Mech to fire at a Target
func (m *Mech) Fire(rangeToTarget int, target weapon.Target) {
	x, y := m.entity.Position()
	for _, w := range m.weapons {
		// Update weapon position before firing
		w.SetPosition(x, y)
		result := w.Fire(rangeToTarget, target)
		if result == false {
			m.notifier.AddMessage("Missed " + target.Name())
		}
	}
}

func (m *Mech) attack(target weapon.Target) {
	if target == nil {
		return
	}
	if target.IsDestroyed() {
		return
	}

	targetX, targetY := target.Position()
	distance := util.CalculateDistance(m.prevX, m.prevY, targetX, targetY)
	m.Fire((int)(distance), target)
	m.game.Log("distance " + strconv.Itoa((int)(distance)))
	m.game.Log("firer (%d,%d), target (%d,%d)", m.prevX, m.prevY, targetX, targetY)
}

// isValidMove checks if a move to the new position is valid
func (m *Mech) isValidMove(newX, newY int) bool {
	// Check game boundaries
	if newX < minCoordinate || newX > maxLevelWidth ||
		newY < minCoordinate || newY > maxLevelHeight {
		if m.game != nil {
			m.game.Log("%s attempted to move out of bounds to (%d,%d)", m.name, newX, newY)
		}
		return false
	}

	// Check for collisions with other entities if we have a level
	if m.level != nil {
		// Check for collisions with other entities
		for _, entity := range m.level.Entities {
			// Skip self and non-physical entities
			if entity == m.entity || entity == nil {
				continue
			}
			
			// Check if entity implements Physical interface
			physical, ok := entity.(tl.Physical)
			if !ok {
				continue
			}

			// Get entity position
			eX, eY := physical.Position()
			
			// If entity is at target position, collision detected
			if eX == newX && eY == newY {
				if m.game != nil {
					m.game.Log("%s attempted to move into occupied position (%d,%d)", m.name, newX, newY)
				}
				return false
			}
		}
	}

	return true
}
