package ai

import (
    "fmt"
    "time"
)

// GameContext represents the current state of the game world
type GameContext struct {
    Player      PlayerInfo       `json:"player"`
    TimeOfDay   string          `json:"time_of_day"`
    Buildings   []BuildingInfo   `json:"buildings"`
    Environment EnvironmentInfo  `json:"environment"`
}

// PlayerInfo contains all relevant information about the player
type PlayerInfo struct {
    Name         string         `json:"name"`
    Occupation   string         `json:"occupation"`
    Money        float64        `json:"money"`
    Health       int            `json:"health"`
    Position     Position       `json:"position"`
    Relationships []Relationship `json:"relationships"`
    Assets       PlayerAssets   `json:"assets"`
}

// Position represents x,y coordinates
type Position struct {
    X int `json:"x"`
    Y int `json:"y"`
}

// Relationship represents a connection with another character
type Relationship struct {
    Name     string `json:"name"`
    Type     string `json:"type"`
    Level    int    `json:"level"`    // 1-10 scale
    Standing string `json:"standing"`  // friendly, neutral, hostile
}

// PlayerAssets represents all assets owned by the player
type PlayerAssets struct {
    Properties []Property `json:"properties"`
    Vehicles   []Vehicle `json:"vehicles"`
    Weapons    []Weapon  `json:"weapons"`
}

// Property represents a real estate property
type Property struct {
    Type      string  `json:"type"`
    Address   string  `json:"address"`
    Value     float64 `json:"value"`
    Condition int     `json:"condition"` // 1-10 scale
}

// Vehicle represents a vehicle owned by the player
type Vehicle struct {
    Type      string  `json:"type"`
    Make      string  `json:"make"`
    Model     string  `json:"model"`
    Year      int     `json:"year"`
    Value     float64 `json:"value"`
    Condition int     `json:"condition"` // 1-10 scale
}

// Weapon represents a weapon in the player's inventory
type Weapon struct {
    Type       string `json:"type"`
    Name       string `json:"name"`
    Damage     int    `json:"damage"`
    Condition  int    `json:"condition"` // 1-10 scale
}

// BuildingInfo represents a building in the game world
type BuildingInfo struct {
    Type        string   `json:"type"`
    Name        string   `json:"name"`
    Position    Position `json:"position"`
    Size        Size     `json:"size"`
    Occupants   int      `json:"occupants"`
    Condition   int      `json:"condition"` // 1-10 scale
    IsHostile   bool     `json:"is_hostile"`
}

// Size represents width and height
type Size struct {
    Width  int `json:"width"`
    Height int `json:"height"`
}

// EnvironmentInfo represents environmental conditions
type EnvironmentInfo struct {
    TimeOfDay    string `json:"time_of_day"`    // morning, afternoon, evening, night
    Visibility   int    `json:"visibility"`      // 1-10 scale
    ThreatLevel  int    `json:"threat_level"`   // 1-10 scale
    ActiveAlerts []string `json:"active_alerts"` // current alerts or warnings
}

// TimeSystem represents the interface for accessing game time
type TimeSystem interface {
    GetCurrentTime() time.Time
}

// ComputerUser represents a computer-controlled character in the game
type ComputerUser struct {
    Name         string
    Age          int
    Nationality  string
    Occupation   string
    PocketMoney  float64
    Properties   []Property
    Cars         []Vehicle
    DailyRoutine DailyRoutine
    Relationships []string
}

// DailyRoutine represents a character's daily schedule
type DailyRoutine struct {
    WakeUpTime string
    SleepTime  string
    Activities []string
}

// NewGameContext creates a new game context with the current state
func NewGameContext(
    player *PlayerInfo,
    buildings []BuildingInfo,
    timeSystem TimeSystem,
) *GameContext {
    // Convert game time to time of day
    timeOfDay := getTimeOfDay(timeSystem.GetCurrentTime())
    
    return &GameContext{
        Player:    *player,
        TimeOfDay: timeOfDay,
        Buildings: buildings,
        Environment: EnvironmentInfo{
            TimeOfDay:    timeOfDay,
            Visibility:   calculateVisibility(timeOfDay),
            ThreatLevel:  calculateThreatLevel(buildings),
            ActiveAlerts: []string{},
        },
    }
}

// getTimeOfDay converts time to a period of the day
func getTimeOfDay(t time.Time) string {
    hour := t.Hour()
    switch {
    case hour >= 5 && hour < 12:
        return "morning"
    case hour >= 12 && hour < 17:
        return "afternoon"
    case hour >= 17 && hour < 21:
        return "evening"
    default:
        return "night"
    }
}

// calculateVisibility returns visibility level based on time of day
func calculateVisibility(timeOfDay string) int {
    switch timeOfDay {
    case "morning", "afternoon":
        return 10
    case "evening":
        return 7
    default: // night
        return 4
    }
}

// calculateThreatLevel analyzes buildings to determine threat level
func calculateThreatLevel(buildings []BuildingInfo) int {
    var hostileCount int
    for _, building := range buildings {
        if building.IsHostile {
            hostileCount++
        }
    }
    // Scale threat level from 1-10 based on hostile buildings
    return min(10, max(1, hostileCount*2))
}

// min returns the smaller of two integers
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// max returns the larger of two integers
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// FormatPrompt creates a natural language prompt from the game context
func (gc *GameContext) FormatPrompt() string {
    return fmt.Sprintf(`Current situation:
Time: %s
Player: %s, a %s with $%.2f
Location: At position (%d, %d)
Health: %d%%

Nearby buildings:
%s

Environment:
- Visibility: %d/10
- Threat Level: %d/10
%s

Assets:
Properties: %d
Vehicles: %d
Weapons: %d

What would be the most strategic course of action for the player?`,
        gc.TimeOfDay,
        gc.Player.Name,
        gc.Player.Occupation,
        gc.Player.Money,
        gc.Player.Position.X,
        gc.Player.Position.Y,
        gc.Player.Health,
        formatBuildings(gc.Buildings),
        gc.Environment.Visibility,
        gc.Environment.ThreatLevel,
        formatAlerts(gc.Environment.ActiveAlerts),
        len(gc.Player.Assets.Properties),
        len(gc.Player.Assets.Vehicles),
        len(gc.Player.Assets.Weapons),
    )
}

// formatBuildings creates a readable list of nearby buildings
func formatBuildings(buildings []BuildingInfo) string {
    if len(buildings) == 0 {
        return "No buildings nearby"
    }
    
    var result string
    for _, b := range buildings {
        status := "Safe"
        if b.IsHostile {
            status = "Hostile"
        }
        result += fmt.Sprintf("- %s (%s) at (%d, %d), Condition: %d/10, Status: %s\n",
            b.Name,
            b.Type,
            b.Position.X,
            b.Position.Y,
            b.Condition,
            status,
        )
    }
    return result
}

// formatAlerts formats active alerts into a readable string
func formatAlerts(alerts []string) string {
    if len(alerts) == 0 {
        return "No active alerts"
    }
    
    var result string
    for _, alert := range alerts {
        result += fmt.Sprintf("- %s\n", alert)
    }
    return result
}
