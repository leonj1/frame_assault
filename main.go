// main.go
package main

import (
	"fmt"
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
	name     string
	color    tl.Attr
	char     rune
	maxCount int
}

var buildingTypes = []BuildingType{
	{"Hospital", tl.ColorRed, 'H', 1},
	{"School", tl.ColorYellow, 'S', 2},
	{"Bank", tl.ColorGreen, 'B', 2},
	{"Grocery", tl.ColorCyan, 'G', 3},
	{"Police", tl.ColorBlue, 'P', 2},
	{"Library", tl.ColorMagenta, 'L', 2},
	{"Mall", tl.ColorWhite, 'M', 2},
	{"Restaurant", tl.ColorRed, 'R', 4},
	{"Theater", tl.ColorYellow, 'T', 2},
	{"Gym", tl.ColorGreen, 'Y', 3},
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

// getValidPatrolPoints generates patrol points that don't overlap with buildings
func getValidPatrolPoints(x, y int, level *tl.BaseLevel) ([][2]int, error) {
	// Try different patrol patterns until we find a valid one
	patterns := []struct {
		dx1, dy1, dx2, dy2 int
	}{
		// Horizontal patrol (left to right)
		{buildingMargin, 1, buildingMargin + buildingSize, 1},
		// Vertical patrol (top to bottom)
		{buildingMargin, 0, buildingMargin, buildingSize},
		// Diagonal patrol
		{buildingMargin, 1, buildingMargin + buildingSize/2, buildingSize/2},
	}

	// Check each pattern for validity
	for _, p := range patterns {
		point1 := [2]int{x + p.dx1, y + p.dy1}
		point2 := [2]int{x + p.dx2, y + p.dy2}

		// Validate points are within bounds
		if !isPointInBounds(point1[0], point1[1]) || !isPointInBounds(point2[0], point2[1]) {
			continue
		}

		// Check for collisions with buildings
		if !hasCollision(point1[0], point1[1], level) && !hasCollision(point2[0], point2[1], level) {
			return [][2]int{point1, point2}, nil
		}
	}

	return nil, fmt.Errorf("no valid patrol points found at position (%d,%d)", x, y)
}

// isPointInBounds checks if a point is within game boundaries
func isPointInBounds(x, y int) bool {
	return x >= minCoordinate && x <= maxLevelWidth &&
		y >= minCoordinate && y <= maxLevelHeight
}

// hasCollision checks if a point collides with any physical entity
func hasCollision(x, y int, level *tl.BaseLevel) bool {
	for _, entity := range level.Entities {
		if entity == nil {
			continue
		}

		physical, ok := entity.(tl.Physical)
		if !ok {
			continue
		}

		// Get entity position and size
		eX, eY := physical.Position()
		if eX == x && eY == y {
			return true
		}
	}
	return false
}

// GenerateEnemyMechs creates a slice of mechs to be used as enemies
func GenerateEnemyMechs(number int, game *tl.Game, level *tl.BaseLevel) []*mech.EnemyMech {
	enemyMechs := make([]*mech.EnemyMech, number)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < number; i++ {
		// Keep trying different positions until we find a valid one
		var strategy movement.Strategy
		var finalX, finalY int

		for attempts := 0; attempts < 10; attempts++ {
			// Random starting position
			x := -15 + r.Intn(30)
			y := -15 + r.Intn(30)

			// Try to get valid patrol points
			patrolPoints, err := getValidPatrolPoints(x, y, level)
			if err != nil {
				if attempts == 9 { // Last attempt, fallback to random walk
					strategy = movement.NewRandomWalkStrategy()
					finalX, finalY = x, y // Use last attempted position
					if game != nil {
						game.Log("Failed to find valid patrol points after %d attempts, using random walk", attempts+1)
					}
				}
				continue
			}

			// Create patrol strategy with valid points
			patrolStrategy, err := movement.NewPatrolStrategy(patrolPoints)
			if err != nil {
				if game != nil {
					game.Log("Failed to create patrol strategy: %v, falling back to random walk", err)
				}
				strategy = movement.NewRandomWalkStrategy()
			} else {
				strategy = patrolStrategy
			}
			finalX, finalY = x, y // Use position where valid patrol points were found
			break
		}

		// If no strategy was created (shouldn't happen due to random walk fallback)
		if strategy == nil {
			strategy = movement.NewRandomWalkStrategy()
		}

		// Create enemy mech using configuration
		config := enemyMechConfigs[i%len(enemyMechConfigs)]
		m := mech.NewEnemyMech(config.name, i, finalX, finalY, tl.ColorRed, config.symbol, strategy)
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
	minCoordinate = 0
	maxLevelWidth = levelWidth - 1
	maxLevelHeight = levelHeight - 1
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

	// Track building counts
	buildingCounts := make(map[string]int)
	for _, bt := range buildingTypes {
		buildingCounts[bt.name] = 0
	}

	// City blocks (buildings)
	buildingIndex := 0
	for x := 0; x < levelWidth-buildingSize; x += avenueSpacing {
		for y := 0; y < levelHeight-buildingSize; y += streetSpacing {
			if x+buildingSize <= levelWidth && y+buildingSize <= levelHeight {
				// Find next available building type
				var buildingType BuildingType
				for attempts := 0; attempts < len(buildingTypes); attempts++ {
					candidateType := buildingTypes[buildingIndex%len(buildingTypes)]
					if buildingCounts[candidateType.name] < candidateType.maxCount {
						buildingType = candidateType
						buildingCounts[buildingType.name]++
						break
					}
					buildingIndex++
				}

				// Skip if no building type available
				if buildingType.name == "" {
					continue
				}

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

// Relationship represents a connection between the user and another person
type Relationship struct {
	PersonName     string
	RelationType   string
	RelationLevel  int // 1-10 scale
}

// Property represents a real estate property owned by the user
type Property struct {
	Address     string
	Type        string
	Value       float64
	YearBought  int
}

// Car represents a vehicle owned by the user
type Car struct {
	Make      string
	Model     string
	Year      int
	Value     float64
}

// DailyRoutine represents the user's daily schedule
type DailyRoutine struct {
	WakeUpTime    string
	SleepTime     string
	Activities    []string
}

// ComputerUser represents a computer user with their personal and professional details
type ComputerUser struct {
	Name                 string
	Age                 int
	Nationality         string
	Occupation          string
	OccupationDesc      string
	DailyRoutine        DailyRoutine
	PersonalityTraits   []string
	ProfInterests       []string
	PersonalInterests   []string
	Skills              []string
	Relationships       []Relationship
	HealthIssues        []string
	PocketMoney         float64
	Properties          []Property
	Cars                []Car
}

// NewComputerUser creates a new instance of ComputerUser with the provided details
func NewComputerUser(name string, age int, nationality string) *ComputerUser {
	return &ComputerUser{
		Name:               name,
		Age:                age,
		Nationality:        nationality,
		PersonalityTraits:  make([]string, 0),
		ProfInterests:      make([]string, 0),
		PersonalInterests:  make([]string, 0),
		Skills:             make([]string, 0),
		Relationships:      make([]Relationship, 0),
		HealthIssues:       make([]string, 0),
		Properties:         make([]Property, 0),
		Cars:              make([]Car, 0),
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
	enemies := GenerateEnemyMechs(8, game, level)
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
