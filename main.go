// main.go
package main

import (
	"math/rand"
	"time"

	"github.com/Ariemeth/frame_assault/display"
	"github.com/Ariemeth/frame_assault/mech"
	"github.com/Ariemeth/frame_assault/mech/weapon"
	tl "github.com/Ariemeth/termloop"
)

// BuildingType represents different types of buildings
type BuildingType struct {
	name  string
	color tl.Attr
	char  rune
}

var buildingTypes = []BuildingType{
	{"Hospital", tl.ColorRed, 'H'},
	{"School", tl.ColorYellow, 'S'},
	{"Bank", tl.ColorGreen, 'B'},
	{"Grocery", tl.ColorCyan, 'G'},
	{"Police", tl.ColorBlue, 'P'},
	{"Library", tl.ColorMagenta, 'L'},
	{"Mall", tl.ColorWhite, 'M'},
	{"Restaurant", tl.ColorRed, 'R'},
	{"Theater", tl.ColorYellow, 'T'},
	{"Gym", tl.ColorGreen, 'Y'},
}

// Building represents a city building with a specific purpose
type Building struct {
	*tl.Entity
	buildingType BuildingType
	width        int
	height       int
}

func NewBuilding(x, y, width, height int, buildingType BuildingType) *Building {
	building := &Building{
		Entity:       tl.NewEntity(x, y, width, height),
		buildingType: buildingType,
		width:        width,
		height:       height,
	}
	return building
}

func (b *Building) Draw(s *tl.Screen) {
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			x, y := b.Position()
			// Draw building outline
			if i == 0 || i == b.width-1 || j == 0 || j == b.height-1 {
				s.RenderCell(x+i, y+j, &tl.Cell{
					Bg: b.buildingType.color,
					Fg: tl.ColorBlack,
					Ch: '█',
				})
			} else if i == b.width/2 && j == b.height/2 {
				// Draw building type identifier in center
				s.RenderCell(x+i, y+j, &tl.Cell{
					Bg: b.buildingType.color,
					Fg: tl.ColorBlack,
					Ch: b.buildingType.char,
				})
			} else {
				// Fill building interior
				s.RenderCell(x+i, y+j, &tl.Cell{
					Bg: b.buildingType.color,
					Fg: b.buildingType.color,
					Ch: ' ',
				})
			}
		}
	}
}

// GenerateEnemyMechs creates a slice of mechs to be used as enemies
func GenerateEnemyMechs(number int) []*mech.Mech {
	enemyMechs := make([]*mech.Mech, 0, number)

	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 1; i <= number; i++ {
		var m *mech.Mech

		chance := i % 8
		x := -15 + r.Intn(30)
		y := -15 + r.Intn(30)

		switch chance {
		case 0:
			m = mech.NewMech("Mech A", i, x, y, tl.ColorRed, rune('A'))
			m.AddWeapon(weapon.CreateRifle())
			break
		case 1:
			m = mech.NewMech("Mech B", i, x, y, tl.ColorRed, rune('B'))
			m.AddWeapon(weapon.CreateRifle())
			break
		case 2:
			m = mech.NewMech("Mech C", i, x, y, tl.ColorRed, rune('C'))
			m.AddWeapon(weapon.CreateShotgun())
			break
		case 3:
			m = mech.NewMech("Mech D", i, x, y, tl.ColorRed, rune('D'))
			m.AddWeapon(weapon.CreateShotgun())
			break
		case 4:
			m = mech.NewMech("Mech E", i, x, y, tl.ColorRed, rune('E'))
			m.AddWeapon(weapon.CreateFist())
			break
		case 5:
			m = mech.NewMech("Mech F", i, x, y, tl.ColorRed, rune('F'))
			m.AddWeapon(weapon.CreateFist())
			break
		case 6:
			m = mech.NewMech("Mech G", i, x, y, tl.ColorRed, rune('G'))
			m.AddWeapon(weapon.CreateSword())
			break
		case 7:
			m = mech.NewMech("Mech H", i, x, y, tl.ColorRed, rune('H'))
			m.AddWeapon(weapon.CreateSword())
			break
		}

		if m != nil {
			enemyMechs = append(enemyMechs, m)
		}
	}

	return enemyMechs
}

// RoadSystem represents a collection of road tiles managed by a single entity
type RoadSystem struct {
	*tl.Entity
	roads map[int]map[int]bool
}

func NewRoadSystem() *RoadSystem {
	return &RoadSystem{
		Entity: tl.NewEntity(0, 0, 1, 1),
		roads:  make(map[int]map[int]bool),
	}
}

func (r *RoadSystem) AddRoad(x, y int) {
	if r.roads[x] == nil {
		r.roads[x] = make(map[int]bool)
	}
	r.roads[x][y] = true
}

func (r *RoadSystem) Draw(s *tl.Screen) {
	for x, yMap := range r.roads {
		for y := range yMap {
			s.RenderCell(x, y, &tl.Cell{
				Bg: tl.ColorBlue,
				Fg: tl.ColorBlue,
				Ch: ' ',
			})
		}
	}
}

const (
	levelWidth     = 60
	levelHeight    = 40
	avenueSpacing  = 12
	streetSpacing  = 8
	buildingWidth  = 5
	buildingHeight = 5
	buildingOffset = 5
)

func createManhattanLayout(level *tl.BaseLevel) {
	roadSystem := NewRoadSystem()

	// Main avenues (vertical roads)
	for x := buildingOffset - 2; x < levelWidth; x += avenueSpacing {
		for y := 0; y < levelHeight; y++ {
			roadSystem.AddRoad(x, y)
			roadSystem.AddRoad(x+1, y)
		}
	}

	// Cross streets (horizontal roads)
	for y := buildingOffset; y < levelHeight; y += streetSpacing {
		for x := 0; x < levelWidth; x++ {
			roadSystem.AddRoad(x, y)
		}
	}

	// Add the road system as a single entity
	level.AddEntity(roadSystem)

	// City blocks (buildings)
	buildingIndex := 0
	for x := 0; x < levelWidth-buildingWidth; x += avenueSpacing {
		for y := 0; y < levelHeight-buildingHeight; y += streetSpacing {
			if x+buildingWidth <= levelWidth && y+buildingHeight <= levelHeight {
				// Cycle through building types
				buildingType := buildingTypes[buildingIndex%len(buildingTypes)]
				building := NewBuilding(
					x+buildingOffset,
					y+1,
					buildingWidth,
					buildingHeight,
					buildingType,
				)
				level.AddEntity(building)
				buildingIndex++
			}
		}
	}
}

func main() {
	//Create the game
	game := tl.NewGame()

	//Create the level for the game
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack, // Dark background for contrast
		Fg: tl.ColorBlack,
		Ch: ' ',
	})

	// Create Manhattan-like layout
	createManhattanLayout(level)

	//Create the notification display
	notification := display.NewNotification(25, 0, 45, 6, level)

	//Create the enemy mechs
	enemies := GenerateEnemyMechs(8)
	for _, enemy := range enemies {
		enemy.AttachGame(game)
		enemy.AttachNotifier(notification)
		enemy.SetLevel(level)
		level.AddEntity(enemy)
	}

	//Create the player's mech
	player := mech.NewPlayerMech("Player", 10, 1, 1, level)
	weapon1 := weapon.CreateRifle()
	player.AddWeapon(weapon1)
	player.SetEnemyList(enemies)
	player.AttachGame(game)
	player.AttachNotifier(notification)
	player.SetLevel(level)
	level.AddEntity(player)

	//Create the players mech status display
	status := display.NewPlayerStatus(0, 0, 20, 13, player, level)

	//Attach the displays the the level
	level.AddEntity(status)
	level.AddEntity(notification)

	//Set the level to be the current game level
	game.Screen().SetLevel(level)

	game.SetDebugOn(false)

	//Start the game engine
	game.Start()
}
