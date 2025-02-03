// main.go
package main

import (
	"math/rand"
	"time"

	"github.com/Ariemeth/frame_assault/display"
	"github.com/Ariemeth/frame_assault/mech"
	"github.com/Ariemeth/frame_assault/mech/movement"
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
	x, y := b.Position()
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
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
func GenerateEnemyMechs(number int) []*mech.EnemyMech {
	enemyMechs := make([]*mech.EnemyMech, 0, number)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create patrol points for different areas of the map
	patrolAreas := [][][2]int{
		{{-15, -15}, {15, -15}, {15, 15}, {-15, 15}}, // Outer patrol
		{{-8, -8}, {8, -8}, {8, 8}, {-8, 8}},         // Inner patrol
		{{-5, 0}, {5, 0}},                             // Horizontal patrol
		{{0, -5}, {0, 5}},                             // Vertical patrol
	}

	for i := 1; i <= number; i++ {
		var m *mech.EnemyMech
		x := -15 + r.Intn(30)
		y := -15 + r.Intn(30)

		// Choose movement strategy
		var strategy movement.Strategy
		if i%2 == 0 {
			strategy = movement.NewRandomWalkStrategy()
		} else {
			patrolPoints := patrolAreas[i%len(patrolAreas)]
			strategy = movement.NewPatrolStrategy(patrolPoints)
		}

		// Create enemy mech with different types
		chance := i % 8
		switch chance {
		case 0:
			m = mech.NewEnemyMech("Mech A", i, x, y, tl.ColorRed, rune('A'), strategy)
			m.AddWeapon(weapon.CreateRifle())
		case 1:
			m = mech.NewEnemyMech("Mech B", i, x, y, tl.ColorRed, rune('B'), strategy)
			m.AddWeapon(weapon.CreateRifle())
		case 2:
			m = mech.NewEnemyMech("Mech C", i, x, y, tl.ColorRed, rune('C'), strategy)
			m.AddWeapon(weapon.CreateShotgun())
		case 3:
			m = mech.NewEnemyMech("Mech D", i, x, y, tl.ColorRed, rune('D'), strategy)
			m.AddWeapon(weapon.CreateShotgun())
		case 4:
			m = mech.NewEnemyMech("Mech E", i, x, y, tl.ColorRed, rune('E'), strategy)
			m.AddWeapon(weapon.CreateSword())
		case 5:
			m = mech.NewEnemyMech("Mech F", i, x, y, tl.ColorRed, rune('F'), strategy)
			m.AddWeapon(weapon.CreateSword())
		case 6:
			m = mech.NewEnemyMech("Mech G", i, x, y, tl.ColorRed, rune('G'), strategy)
			m.AddWeapon(weapon.CreateFist())
		case 7:
			m = mech.NewEnemyMech("Mech H", i, x, y, tl.ColorRed, rune('H'), strategy)
			m.AddWeapon(weapon.CreateFist())
		}

		enemyMechs = append(enemyMechs, m)
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
	streetSpacing  = 6
	buildingMargin = 2
	buildingSize   = 4
	gameFPS       = 10 // Run at 10 FPS for smoother animation while keeping slow movement
)

func createManhattanLayout(level *tl.BaseLevel) {
	roadSystem := NewRoadSystem()

	// Main avenues (vertical roads)
	for x := buildingMargin - 2; x < levelWidth; x += avenueSpacing {
		for y := 0; y < levelHeight; y++ {
			roadSystem.AddRoad(x, y)
			roadSystem.AddRoad(x+1, y)
		}
	}

	// Cross streets (horizontal roads)
	for y := buildingMargin; y < levelHeight; y += streetSpacing {
		for x := 0; x < levelWidth; x++ {
			roadSystem.AddRoad(x, y)
		}
	}

	// Add the road system as a single entity
	level.AddEntity(roadSystem)

	// City blocks (buildings)
	buildingIndex := 0
	for x := 0; x < levelWidth-buildingSize; x += avenueSpacing {
		for y := 0; y < levelHeight-buildingSize; y += streetSpacing {
			if x+buildingSize <= levelWidth && y+buildingSize <= levelHeight {
				// Cycle through building types
				buildingType := buildingTypes[buildingIndex%len(buildingTypes)]
				building := NewBuilding(
					x+buildingMargin,
					y+1,
					buildingSize,
					buildingSize,
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
	game.Screen().SetFps(gameFPS)

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
	enemyMechs := make([]*mech.Mech, len(enemies))
	for i, enemy := range enemies {
		enemy.AttachGame(game)
		enemy.AttachNotifier(notification)
		enemy.SetLevel(level)
		level.AddEntity(enemy)
		enemyMechs[i] = enemy.Mech
	}

	//Create the player mech
	player := mech.NewPlayerMech("Player", 10, 0, 0, level)
	player.AttachGame(game)
	player.AttachNotifier(notification)
	player.SetEnemyList(enemyMechs)
	player.AddWeapon(weapon.CreateRifle())
	level.AddEntity(player)

	//Create the player status display
	playerStatus := display.NewPlayerStatus(0, 0, 25, 6, player, level)
	level.AddEntity(playerStatus)

	//Create the notification display
	level.AddEntity(notification)

	//Set the level
	game.Screen().SetLevel(level)

	//Start the game
	game.Start()
}
