package projectile

import (
	"math"
	"time"

	tl "github.com/Ariemeth/termloop"
)

// Bullet represents a projectile fired from a weapon
type Bullet struct {
	*tl.Entity
	x, y             float64 // Current position as float for smooth movement
	targetX, targetY int     // Target position
	dx, dy           float64 // Direction vector
	speed            float64
	symbol           rune
	color            tl.Attr
	level            *tl.BaseLevel
	lastMove         time.Time
	moveDelay        time.Duration
	trail            [][2]float64 // Trail positions
	trailLength      int
}

// NewBullet creates a new bullet entity
func NewBullet(startX, startY, targetX, targetY int, level *tl.BaseLevel) *Bullet {
	bullet := &Bullet{
		Entity:      tl.NewEntity(startX, startY, 1, 1),
		x:           float64(startX),
		y:           float64(startY),
		targetX:     targetX,
		targetY:     targetY,
		speed:       1.0,
		symbol:      '*',
		color:       tl.ColorYellow | tl.AttrBold,
		level:       level,
		lastMove:    time.Now(),
		moveDelay:   time.Millisecond * 100,
		trail:       make([][2]float64, 0),
		trailLength: 3, // Number of trailing bullets
	}

	// Calculate direction vector
	dx := float64(targetX) - bullet.x
	dy := float64(targetY) - bullet.y
	length := math.Sqrt(dx*dx + dy*dy)
	if length != 0 {
		bullet.dx = dx / length
		bullet.dy = dy / length
	}

	return bullet
}

// Draw implements the Draw method of the Drawable interface
func (b *Bullet) Draw(screen *tl.Screen) {
	// Draw trail
	for _, pos := range b.trail {
		screenX := int(math.Round(pos[0]))
		screenY := int(math.Round(pos[1]))
		// Make trail bullets slightly dimmer
		screen.RenderCell(screenX, screenY, &tl.Cell{
			Fg: b.color & ^tl.AttrBold,
			Ch: b.symbol,
		})
	}

	// Draw current bullet position
	screenX := int(math.Round(b.x))
	screenY := int(math.Round(b.y))
	screen.RenderCell(screenX, screenY, &tl.Cell{
		Fg: b.color,
		Ch: b.symbol,
	})
}

// Tick implements the Tick method of the Drawable interface
func (b *Bullet) Tick(event tl.Event) {
	// Only move if enough time has passed
	if time.Since(b.lastMove) < b.moveDelay {
		return
	}

	// Add current position to trail
	b.trail = append(b.trail, [2]float64{b.x, b.y})
	if len(b.trail) > b.trailLength {
		b.trail = b.trail[1:]
	}

	// Update position using floating-point coordinates
	b.x += b.dx * b.speed
	b.y += b.dy * b.speed

	// Convert to screen coordinates
	screenX := int(math.Round(b.x))
	screenY := int(math.Round(b.y))

	// Check if bullet reached target
	if math.Abs(float64(b.targetX)-b.x) < 0.5 && math.Abs(float64(b.targetY)-b.y) < 0.5 {
		b.level.RemoveEntity(b)
		return
	}

	// Update entity position
	b.SetPosition(screenX, screenY)
	b.lastMove = time.Now()
}
