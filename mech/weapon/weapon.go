package weapon

import (
	"math/rand"
	"time"

	"github.com/Ariemeth/frame_assault/projectile"
	tl "github.com/Ariemeth/termloop"
)

// Weapon is weapon with specific characteristics
type Weapon struct {
	maxRange, damage int
	name             string
	hitRate          float64
	level            *tl.BaseLevel
	sourceX, sourceY int // Position of the weapon holder
}

// Target is an interface used by objects that can be hit and take damage
type Target interface {
	// Hit is called when an object is hit and the amount of damage to be done.
	Hit(int)
	// Name should return the name of the target.
	Name() string
	// IsDestroyed should return true is the target is destroyed, false otherwise.
	IsDestroyed() bool
	// Position should return the x,y location of the target.
	Position() (int, int)
}

// Create creates a new Weapon.
func Create(maxRange int, damage int, name string,
	hitRate float64) Weapon {

	return Weapon{maxRange: maxRange, damage: damage, name: name,
		hitRate: hitRate}
}

// SetLevel sets the game level reference for creating bullets
func (weapon *Weapon) SetLevel(level *tl.BaseLevel) {
	weapon.level = level
}

// SetPosition sets the current position of the weapon holder
func (weapon *Weapon) SetPosition(x, y int) {
	weapon.sourceX = x
	weapon.sourceY = y
}

// Name returns the name of the weapon
func (weapon Weapon) Name() string {
	return weapon.name
}

// Range returns the range of the weapon
func (weapon Weapon) Range() int {
	return weapon.maxRange
}

// Damage returns the damage of the weapon
func (weapon Weapon) Damage() int {
	return weapon.damage
}

// Accuracy returns the accuracy of the weapon
func (weapon Weapon) Accuracy() float64 {
	return weapon.hitRate
}

// Fire is used by an object to fire at a Target.
// Requires the range to the Target and the Target.
// Returns true if the target is hit or false if the target is missed.
func (weapon Weapon) Fire(rangeToTarget int, target Target) bool {
	if rangeToTarget <= weapon.maxRange {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		chance := r.Float64()

		// Create bullet regardless of hit/miss
		if weapon.level != nil {
			targetX, targetY := target.Position()
			bullet := projectile.NewBullet(weapon.sourceX, weapon.sourceY, targetX, targetY, weapon.level)
			weapon.level.AddEntity(bullet)
		}

		if chance <= weapon.Accuracy() {
			target.Hit(weapon.damage)
			return true
		}
	}
	return false
}
