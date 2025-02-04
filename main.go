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
	
	// Draw building outline and fill
	for i := 0; i < b.width; i++ {
		for j := 0; j < b.height; j++ {
			// Draw building outline
			if i == 0 || i == b.width-1 || j == 0 || j == b.height-1 {
				s.RenderCell(x+i, y+j, &tl.Cell{
					Bg: b.buildingType.color,
					Fg: tl.ColorBlack,
					Ch: '█',
				})
			} else {
				// Fill building interior
				s.RenderCell(x+i, y+j, &tl.Cell{
					Bg: b.buildingType.color,
					Fg: tl.ColorBlack,
					Ch: ' ',
				})
			}
		}
	}
	
	// Draw building name in the center
	name := b.buildingType.name
	startX := x + (b.width-len(name))/2
	startY := y + b.height/2
	
	for i, ch := range name {
		if startX+i < x+b.width-1 { // Ensure we don't write outside building bounds
			s.RenderCell(startX+i, startY, &tl.Cell{
				Bg: b.buildingType.color,
				Fg: tl.ColorBlack,
				Ch: ch,
			})
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
		{buildingMargin, 1, buildingMargin + buildingWidth, 1},
		// Vertical patrol (top to bottom)
		{buildingMargin, 0, buildingMargin, buildingHeight},
		// Diagonal patrol
		{buildingMargin, 1, buildingMargin + buildingWidth/2, buildingHeight/2},
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

func (r *RoadSystem) HasRoad(x, y int) bool {
	if yMap, exists := r.roads[x]; exists {
		return yMap[y]
	}
	return false
}

func (r *RoadSystem) HasRoadInArea(x, y, width, height int) bool {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			if r.HasRoad(i, j) {
				return true
			}
		}
	}
	return false
}

const (
	levelWidth     = 100
	levelHeight    = 60
	avenueSpacing  = 20
	streetSpacing  = 12
	buildingMargin = 2
	buildingWidth  = 8
	buildingHeight = 6
	gameFPS       = 10
	minCoordinate = 0
	maxLevelWidth = levelWidth - 1
	maxLevelHeight = levelHeight - 1
    
    // Time constants
    realSecondsPerGameDay = 180.0  // 3 minutes real time = 24 hours game time
    gameHoursPerRealSecond = 24.0 / realSecondsPerGameDay
    gameHoursPerFrame = gameHoursPerRealSecond / gameFPS
    timeDisplayX = 1
    timeDisplayY = 1
)

// TimeSystemInterface defines the interface for time systems
type TimeSystemInterface interface {
	Tick(event tl.Event)
}

// TimeSystem handles the game's time progression
type TimeSystem struct {
	*tl.Entity
	gameHours    float64
	frameCounter int
}

// NewTimeSystem creates a new time system starting at 6:00 AM
func NewTimeSystem(level *tl.BaseLevel) *TimeSystem {
	ts := &TimeSystem{
		Entity:     tl.NewEntity(timeDisplayX, timeDisplayY, 20, 1),
		gameHours:  6.0, // Start at 6 AM
	}
	return ts
}

// FormatGameTime converts game hours to a 12-hour time string
func (ts *TimeSystem) FormatGameTime() string {
	hours := int(ts.gameHours) % 24
	minutes := int((ts.gameHours - float64(int(ts.gameHours))) * 60)
	period := "AM"
	
	if hours >= 12 {
		period = "PM"
		if hours > 12 {
			hours -= 12
		}
	}
	if hours == 0 {
		hours = 12
	}
	
	return fmt.Sprintf("Time: %02d:%02d %s", hours, minutes, period)
}

// Tick updates the game time
func (ts *TimeSystem) Tick(event tl.Event) {
	ts.frameCounter++
	ts.gameHours += gameHoursPerFrame
	if ts.gameHours >= 24.0 {
		ts.gameHours -= 24.0
	}
}

// getSafeSpawnPosition returns a position that is not on a road or building
func getSafeSpawnPosition() (x, y int) {
	// Position player in the middle of a block between roads
	// Add buildingMargin to avoid spawning too close to buildings
	x = buildingMargin + avenueSpacing/2
	y = buildingMargin + streetSpacing/2
	return x, y
}

// createRoadSystem creates and returns a road system with vertical and horizontal roads
func createRoadSystem() *RoadSystem {
	roadSystem := NewRoadSystem()
	createVerticalRoads(roadSystem)
	createHorizontalRoads(roadSystem)
	return roadSystem
}

// createVerticalRoads adds vertical roads (avenues) to the road system
func createVerticalRoads(roadSystem *RoadSystem) {
	for x := buildingMargin - 2; x < levelWidth; x += avenueSpacing {
		for y := 0; y < levelHeight; y++ {
			roadSystem.AddRoad(x, y)
			roadSystem.AddRoad(x+1, y)
		}
	}
}

// createHorizontalRoads adds horizontal roads (streets) to the road system
func createHorizontalRoads(roadSystem *RoadSystem) {
	for y := buildingMargin; y < levelHeight; y += streetSpacing {
		for x := 0; x < levelWidth; x++ {
			roadSystem.AddRoad(x, y)
		}
	}
}

// initBuildingCounts initializes a map to track building counts
func initBuildingCounts() map[string]int {
	counts := make(map[string]int)
	for _, bt := range buildingTypes {
		counts[bt.name] = 0
	}
	return counts
}

// getValidBuildingPositions returns a list of valid positions for building placement
func getValidBuildingPositions(roadSystem *RoadSystem) [][2]int {
	valid := make([][2]int, 0)
	for x := buildingMargin; x < levelWidth-buildingWidth; x += avenueSpacing {
		for y := buildingMargin + 1; y < levelHeight-buildingHeight; y += streetSpacing {
			// Check if the entire building area is free of roads
			if !roadSystem.HasRoadInArea(x, y, buildingWidth, buildingHeight) {
				valid = append(valid, [2]int{x, y})
			}
		}
	}
	return valid
}

// tryPlaceBuilding attempts to place a building at the given location
func tryPlaceBuilding(x, y int, buildingCounts map[string]int, level *tl.BaseLevel) bool {
	for tries := 0; tries < len(buildingTypes)*2; tries++ {
		buildingType := buildingTypes[rand.Intn(len(buildingTypes))]
		if buildingCounts[buildingType.name] < buildingType.maxCount {
			building := NewBuilding(x, y, buildingWidth, buildingHeight, buildingType)
			level.AddEntity(building)
			buildingCounts[buildingType.name]++
			return true
		}
	}
	return false
}

// placeBuildings places buildings in valid positions
func placeBuildings(roadSystem *RoadSystem, buildingCounts map[string]int, level *tl.BaseLevel) {
	validPositions := getValidBuildingPositions(roadSystem)
	for _, pos := range validPositions {
		tryPlaceBuilding(pos[0], pos[1], buildingCounts, level)
	}
}

// createManhattanLayout creates the city layout with roads and buildings
func createManhattanLayout(level *tl.BaseLevel) {
	roadSystem := createRoadSystem()
	level.AddEntity(roadSystem)
	
	buildingCounts := initBuildingCounts()
	placeBuildings(roadSystem, buildingCounts, level)
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

const (
	lowIncomeUsers    = 0.6
	middleIncomeUsers = 0.3
	highIncomeUsers   = 0.1
)

// IncomeLevel represents different income levels for computer users
type IncomeLevel int

const (
	LowIncome IncomeLevel = iota
	MiddleIncome
	HighIncome
)

// GenerateComputerUsers creates a slice of computer users with varying income levels
func GenerateComputerUsers(number int) []*ComputerUser {
	users := make([]*ComputerUser, number)
	
	// Calculate number of users per income level
	lowCount := int(float64(number) * lowIncomeUsers)
	middleCount := int(float64(number) * middleIncomeUsers)
	highCount := number - lowCount - middleCount
	
	currentIndex := 0
	
	// Generate low income users
	for i := 0; i < lowCount; i++ {
		users[currentIndex] = generateUserByIncomeLevel(LowIncome)
		currentIndex++
	}
	
	// Generate middle income users
	for i := 0; i < middleCount; i++ {
		users[currentIndex] = generateUserByIncomeLevel(MiddleIncome)
		currentIndex++
	}
	
	// Generate high income users
	for i := 0; i < highCount; i++ {
		users[currentIndex] = generateUserByIncomeLevel(HighIncome)
		currentIndex++
	}
	
	return users
}

// generateUserByIncomeLevel creates a single computer user based on their income level
func generateUserByIncomeLevel(level IncomeLevel) *ComputerUser {
	nationalities := []string{"American", "Canadian", "British", "German", "Japanese", "Australian"}
	occupations := map[IncomeLevel][]string{
		LowIncome: {"Retail Worker", "Server", "Delivery Driver", "Security Guard"},
		MiddleIncome: {"Teacher", "Nurse", "Office Manager", "Sales Representative"},
		HighIncome: {"Software Engineer", "Doctor", "Lawyer", "Business Executive"},
	}
	
	// Random name generation (simple version)
	firstNames := []string{"John", "Jane", "Mike", "Sarah", "David", "Emma"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia"}
	
	rand.Seed(time.Now().UnixNano())
	name := firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
	age := 20 + rand.Intn(45) // Ages between 20 and 65
	nationality := nationalities[rand.Intn(len(nationalities))]
	
	user := NewComputerUser(name, age, nationality)
	
	// Set occupation based on income level
	possibleOccupations := occupations[level]
	user.Occupation = possibleOccupations[rand.Intn(len(possibleOccupations))]
	
	// Set income level specific attributes
	switch level {
	case LowIncome:
		user.PocketMoney = 500 + float64(rand.Intn(1500))
		if rand.Float64() < 0.3 { // 30% chance to have a car
			user.Cars = append(user.Cars, Car{
				Make:  "Toyota",
				Model: "Corolla",
				Year:  2010 + rand.Intn(5),
				Value: 5000 + float64(rand.Intn(3000)),
			})
		}
	
	case MiddleIncome:
		user.PocketMoney = 3000 + float64(rand.Intn(4000))
		user.Cars = append(user.Cars, Car{
			Make:  "Honda",
			Model: "Accord",
			Year:  2015 + rand.Intn(5),
			Value: 15000 + float64(rand.Intn(10000)),
		})
		if rand.Float64() < 0.4 { // 40% chance to have a property
			user.Properties = append(user.Properties, Property{
				Address:    "123 Suburban St",
				Type:      "House",
				Value:     250000 + float64(rand.Intn(150000)),
				YearBought: 2015 + rand.Intn(8),
			})
		}
	
	case HighIncome:
		user.PocketMoney = 10000 + float64(rand.Intn(40000))
		// Multiple cars
		cars := []Car{
			{
				Make:  "BMW",
				Model: "5 Series",
				Year:  2020 + rand.Intn(4),
				Value: 50000 + float64(rand.Intn(30000)),
			},
			{
				Make:  "Tesla",
				Model: "Model S",
				Year:  2021 + rand.Intn(3),
				Value: 80000 + float64(rand.Intn(40000)),
			},
		}
		user.Cars = cars
		
		// Multiple properties
		properties := []Property{
			{
				Address:    "456 Luxury Ave",
				Type:      "House",
				Value:     800000 + float64(rand.Intn(500000)),
				YearBought: 2018 + rand.Intn(5),
			},
			{
				Address:    "789 Investment St",
				Type:      "Rental Property",
				Value:     400000 + float64(rand.Intn(200000)),
				YearBought: 2016 + rand.Intn(7),
			},
		}
		user.Properties = properties
	}
	
	// Set daily routine
	user.DailyRoutine = DailyRoutine{
		WakeUpTime: "07:00",
		SleepTime:  "23:00",
		Activities: []string{"Work", "Exercise", "Leisure"},
	}
	
	return user
}

// ComputerUserEntity represents a visual entity for a computer user in the game
type ComputerUserEntity struct {
	*tl.Entity
	user *ComputerUser
	symbol rune
	color tl.Attr
}

// NewComputerUserEntity creates a new computer user entity for rendering
func NewComputerUserEntity(user *ComputerUser, x, y int) *ComputerUserEntity {
	// Different symbols and colors based on income level
	var symbol rune
	var color tl.Attr
	
	// Determine pocket money to set income level
	switch {
	case user.PocketMoney >= 10000: // High income
		symbol = '⚫' // Rich user symbol
		color = tl.ColorGreen
	case user.PocketMoney >= 3000: // Middle income
		symbol = '◉' // Middle class symbol
		color = tl.ColorYellow
	default: // Low income
		symbol = '○' // Low income symbol
		color = tl.ColorRed
	}
	
	return &ComputerUserEntity{
		Entity: tl.NewEntity(x, y, 1, 1),
		user:   user,
		symbol: symbol,
		color:  color,
	}
}

// Draw implements the termloop.Drawable interface
func (c *ComputerUserEntity) Draw(screen *tl.Screen) {
	x, y := c.Position()
	screen.RenderCell(x, y, &tl.Cell{
		Fg: c.color,
		Ch: c.symbol,
	})
}

// Tick implements the termloop.Drawable interface
func (c *ComputerUserEntity) Tick(event tl.Event) {
	// For now, users stay in place
	// TODO: Implement movement patterns based on daily routine
}

// Collide implements termloop.Physical interface
func (c *ComputerUserEntity) Collide(collision tl.Physical) {
	// Handle collisions if needed
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
    
    // Create and add time system
    timeSystem := NewTimeSystem(level)
    level.AddEntity(timeSystem)
    
    // Generate and place computer users
    users := GenerateComputerUsers(8)
    
    // Place users around the map at predefined positions
    userPositions := []struct{ x, y int }{
        {5, 5}, {10, 5}, {15, 5}, {20, 5},  // Top row
        {5, 10}, {10, 10}, {15, 10}, {20, 10}, // Bottom row
    }
    
    for i, user := range users {
        pos := userPositions[i]
        userEntity := NewComputerUserEntity(user, pos.x, pos.y)
        level.AddEntity(userEntity)
    }
    
    //Create the enemy mechs
    enemies := GenerateEnemyMechs(8, game, level)
    enemyMechs := make([]*mech.Mech, len(enemies))
    
    for i, enemy := range enemies {
        enemy.SetLevel(level)
        enemy.AttachNotifier(notification)
        level.AddEntity(enemy)
        enemyMechs[i] = enemy.Mech
    }

    // Create the player mech
    spawnX := 1
    spawnY := 1
    player := mech.NewPlayerMech("Player", 10, spawnX, spawnY, level)
    player.AttachGame(game)
    player.AttachNotifier(notification)
    player.SetEnemyList(enemyMechs)
    level.AddEntity(player)
    player.AddWeapon(weapon.CreateRifle())
    
    // Create player status display with time system
    playerStatus := display.NewPlayer(0, 0, player, timeSystem, level)
    level.AddEntity(playerStatus)

	//Create the player status display
	level.AddEntity(notification)

	//Set the level
	game.Screen().SetLevel(level)

	//Start the game
	game.Start()
}
