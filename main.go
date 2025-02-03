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

// mechConfig defines the configuration for creating an enemy mech
type mechConfig struct {
	name     string
	symbol   rune
	weapon   func() weapon.Weapon
}

// enemyMechConfigs defines the available enemy mech configurations
var enemyMechConfigs = []mechConfig{
	{"Mech A", 'A', weapon.CreateRifle},
	{"Mech B", 'B', weapon.CreateRifle},
	{"Mech C", 'C', weapon.CreateShotgun},
	{"Mech D", 'D', weapon.CreateShotgun},
	{"Mech E", 'E', weapon.CreateSword},
	{"Mech F", 'F', weapon.CreateSword},
	{"Mech G", 'G', weapon.CreateFist},
	{"Mech H", 'H', weapon.CreateFist},
}

// GenerateEnemyMechs creates a slice of mechs to be used as enemies
func GenerateEnemyMechs(number int, game *tl.Game) []*mech.EnemyMech {
	enemyMechs := make([]*mech.EnemyMech, number)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < number; i++ {
		// Random starting position
		x := -15 + r.Intn(30)
		y := -15 + r.Intn(30)

		// Create patrol points for the enemy
		patrolPoints := [][2]int{
			{x + buildingMargin, y + 1},
			{x + buildingMargin + buildingSize, y + 1},
		}

		// Create movement strategy
		var strategy movement.Strategy
		patrolStrategy, err := movement.NewPatrolStrategy(patrolPoints)
		if err != nil {
			// If patrol strategy fails, fallback to random walk
			strategy = movement.NewRandomWalkStrategy()
			if game != nil {
				game.Log("Failed to create patrol strategy: %v, falling back to random walk", err)
			}
		} else {
			strategy = patrolStrategy
		}

		// Create enemy mech using configuration
		config := enemyMechConfigs[i%len(enemyMechConfigs)]
		m := mech.NewEnemyMech(config.name, i, x, y, tl.ColorRed, config.symbol, strategy)
		m.AddWeapon(config.weapon())
		m.AttachGame(game)
		enemyMechs[i] = m
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

// getSafeSpawnPosition returns a position that is not on a road or building
func getSafeSpawnPosition() (x, y int) {
	// Position player in the middle of a block between roads
	// Add buildingMargin to avoid spawning too close to buildings
	x = buildingMargin + avenueSpacing/2
	y = buildingMargin + streetSpacing/2
	return x, y
}

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
	enemies := GenerateEnemyMechs(8, game)
	enemyMechs := make([]*mech.Mech, len(enemies))
	for i, enemy := range enemies {
		enemy.SetLevel(level)
		enemy.AttachNotifier(notification)
		level.AddEntity(enemy)
		enemyMechs[i] = enemy.Mech
	}

	//Create the player mech
	spawnX, spawnY := getSafeSpawnPosition()
	player := mech.NewPlayerMech("Player", 10, spawnX, spawnY, level)
	player.AttachGame(game)
	player.AttachNotifier(notification)
	player.SetEnemyList(enemyMechs)
	level.AddEntity(player)
	player.AddWeapon(weapon.CreateRifle())

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
