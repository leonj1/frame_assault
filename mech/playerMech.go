package mech

import (
	"strings"

	tl "github.com/Ariemeth/termloop"
)

//PlayerMech represents a player controlled mech
type PlayerMech struct {
	Mech
	level   *tl.BaseLevel
	enemies []*Mech
}

// NewPlayerMech is used to create a new instance of a mech with default structure.
func NewPlayerMech(name string, maxStructure, x, y int, level *tl.BaseLevel) *PlayerMech {
	newMech := NewMech(name, maxStructure, x, y, tl.ColorRed, 'M')
	newMech.SetLevel(level)

	newPlayerMech := PlayerMech{
		Mech:  *newMech,
		level: level,
	}

	return &newPlayerMech
}

//SetEnemyList sets the list of enemies the player can interact
func (pMech *PlayerMech) SetEnemyList(enemies []*Mech) {
	pMech.enemies = enemies
}

// Tick is called to process 1 tick of actions based on the
// type of event.
func (pMech *PlayerMech) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		pMech.prevX, pMech.prevY = pMech.entity.Position()

		//quick fix to handle keys
		switch event.Ch {
		case 'A':
		case 'a':
			pMech.attack("A")
			break
		case 'B':
		case 'b':
			pMech.attack("B")
			break
		case 'C':
		case 'c':
			pMech.attack("C")
			break
		case 'D':
		case 'd':
			pMech.attack("D")
			break
		case 'E':
		case 'e':
			pMech.attack("E")
			break
		case 'F':
		case 'f':
			pMech.attack("F")
			break
		case 'G':
		case 'g':
			pMech.attack("G")
			break
		case 'H':
		case 'h':
			pMech.attack("H")
			break
		}

		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			pMech.entity.SetPosition(pMech.prevX+1, pMech.prevY)
			break
		case tl.KeyArrowLeft:
			pMech.entity.SetPosition(pMech.prevX-1, pMech.prevY)
			break
		case tl.KeyArrowUp:
			pMech.entity.SetPosition(pMech.prevX, pMech.prevY-1)
			break
		case tl.KeyArrowDown:
			pMech.entity.SetPosition(pMech.prevX, pMech.prevY+1)
			break
		}
	}
}

// Draw passes the draw call to entity.
func (pMech *PlayerMech) Draw(screen *tl.Screen) {
	screenWidth, screenHeight := screen.Size()
	x, y := pMech.entity.Position()
	pMech.level.SetOffset(screenWidth/2-x, screenHeight/2-y)
	pMech.entity.Draw(screen)
}

func (pMech *PlayerMech) getTargetEnemy(name string) *Mech {
	for i, mech := range pMech.enemies {
		if strings.HasSuffix(mech.Name(), name) {
			pMech.game.Log("enemy found: %s", mech.Name())
			return pMech.enemies[i]
		}
	}
	return nil
}

func (pMech *PlayerMech) attack(name string) {
	target := pMech.getTargetEnemy(name)
	pMech.Mech.attack(target)
}
