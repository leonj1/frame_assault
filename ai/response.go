package ai

import (
    "encoding/json"
    "fmt"
)

// EmotionalState represents the emotional state of a character
type EmotionalState string

const (
    EmotionHappy    EmotionalState = "happy"
    EmotionTired    EmotionalState = "tired"
    EmotionStressed EmotionalState = "stressed"
    EmotionSad      EmotionalState = "sad"
    EmotionAfraid   EmotionalState = "afraid"
    EmotionWorried  EmotionalState = "worried"
    EmotionCalm     EmotionalState = "calm"
    EmotionAngry    EmotionalState = "angry"
    EmotionPanic    EmotionalState = "panic"
)

// ActionPriority indicates how urgent an action is
type ActionPriority string

const (
    PriorityImmediate ActionPriority = "immediate" // Must be done right now
    PriorityHigh      ActionPriority = "high"      // Should be done very soon
    PriorityMedium    ActionPriority = "medium"    // Should be done when convenient
    PriorityLow       ActionPriority = "low"       // Can be done later
)

// ActionType categorizes different types of actions
type ActionType string

const (
    ActionMove     ActionType = "move"      // Movement related actions
    ActionCombat   ActionType = "combat"    // Combat related actions
    ActionSocial   ActionType = "social"    // Social interactions
    ActionWork     ActionType = "work"      // Work related activities
    ActionRest     ActionType = "rest"      // Rest or recovery actions
    ActionDefense  ActionType = "defense"   // Defensive actions
    ActionFlee     ActionType = "flee"      // Escape actions
    ActionExplore  ActionType = "explore"   // Exploration actions
)

// NPCAction represents a specific action the NPC will take
type NPCAction struct {
    Type        ActionType     `json:"type"`
    Priority    ActionPriority `json:"priority"`
    Description string         `json:"description"`
    Target      *Position     `json:"target,omitempty"`     // Where the action will take place
    Duration    int           `json:"duration,omitempty"`   // Estimated time in game minutes
}

// NPCIntent represents the character's current goals and motivations
type NPCIntent struct {
    PrimaryGoal     string   `json:"primary_goal"`      // Main objective
    SecondaryGoals  []string `json:"secondary_goals"`   // Other objectives
    Concerns        []string `json:"concerns"`          // Current worries or issues
    TargetLocation  *Position `json:"target_location,omitempty"` // Where they ultimately want to go
}

// EmotionalResponse represents the character's emotional state and its causes
type EmotionalResponse struct {
    PrimaryEmotion   EmotionalState  `json:"primary_emotion"`
    SecondaryEmotion EmotionalState  `json:"secondary_emotion,omitempty"`
    Intensity        int             `json:"intensity"`        // 1-10 scale
    Causes           []string        `json:"causes"`          // Reasons for emotional state
    PhysicalSigns    []string        `json:"physical_signs"`  // Observable physical manifestations
}

// NPCResponse represents the complete AI response for an NPC's next actions
type NPCResponse struct {
    NextActions []NPCAction       `json:"next_actions"`   // Ordered list of actions to take
    Intent      NPCIntent         `json:"intent"`         // Current goals and motivations
    Emotional   EmotionalResponse `json:"emotional"`      // Emotional state and its causes
    Thoughts    []string          `json:"thoughts"`       // Internal monologue or reasoning
}

// ParseOllamaResponse attempts to parse the Ollama response into a structured NPCResponse
func ParseOllamaResponse(response string) (*NPCResponse, error) {
    var npcResponse NPCResponse
    if err := json.Unmarshal([]byte(response), &npcResponse); err != nil {
        return nil, fmt.Errorf("failed to parse Ollama response: %v", err)
    }
    return &npcResponse, nil
}

// FormatPrompt creates a prompt template for generating NPC responses
func FormatNPCPrompt(context *GameContext, npcInfo *ComputerUser) string {
    return fmt.Sprintf(`Based on the following situation, generate a JSON response describing the NPC's next actions, intent, and emotional state.
The response should follow this structure:
{
    "next_actions": [
        {
            "type": "move|combat|social|work|rest|defense|flee|explore",
            "priority": "immediate|high|medium|low",
            "description": "Detailed description of the action",
            "target": {"x": 0, "y": 0},
            "duration": 30
        }
    ],
    "intent": {
        "primary_goal": "Main objective",
        "secondary_goals": ["Other objectives"],
        "concerns": ["Current worries"],
        "target_location": {"x": 0, "y": 0}
    },
    "emotional": {
        "primary_emotion": "happy|tired|stressed|sad|afraid|worried|calm|angry|panic",
        "secondary_emotion": "happy|tired|stressed|sad|afraid|worried|calm|angry|panic",
        "intensity": 7,
        "causes": ["Reasons for emotional state"],
        "physical_signs": ["Observable physical manifestations"]
    },
    "thoughts": ["Internal monologue", "Reasoning"]
}

Current situation:
Time: %s
NPC: %s, a %s with $%.2f
Location: At position (%d, %d)
Daily Schedule: Wakes up at %s, sleeps at %s
Activities: %v

Environment:
Visibility: %d/10
Threat Level: %d/10
%s

Assets:
Properties: %d
Vehicles: %d
Current relationships: %v

Please respond with a valid JSON object following the structure above.`,
        context.TimeOfDay,
        npcInfo.Name,
        npcInfo.Occupation,
        npcInfo.PocketMoney,
        // TODO: Get actual position from entity
        0, 0,
        npcInfo.DailyRoutine.WakeUpTime,
        npcInfo.DailyRoutine.SleepTime,
        npcInfo.DailyRoutine.Activities,
        context.Environment.Visibility,
        context.Environment.ThreatLevel,
        formatAlerts(context.Environment.ActiveAlerts),
        len(npcInfo.Properties),
        len(npcInfo.Cars),
        npcInfo.Relationships,
    )
}

// ValidateResponse checks if the NPCResponse is valid and complete
func (r *NPCResponse) ValidateResponse() error {
    if len(r.NextActions) == 0 {
        return fmt.Errorf("no actions specified")
    }
    
    if r.Intent.PrimaryGoal == "" {
        return fmt.Errorf("no primary goal specified")
    }
    
    if r.Emotional.PrimaryEmotion == "" {
        return fmt.Errorf("no primary emotion specified")
    }
    
    if r.Emotional.Intensity < 1 || r.Emotional.Intensity > 10 {
        return fmt.Errorf("emotional intensity must be between 1 and 10")
    }
    
    return nil
}

// IsUrgentAction checks if any of the next actions are immediate priority
func (r *NPCResponse) IsUrgentAction() bool {
    for _, action := range r.NextActions {
        if action.Priority == PriorityImmediate {
            return true
        }
    }
    return false
}

// GetNextAction returns the next action to execute, or nil if no actions
func (r *NPCResponse) GetNextAction() *NPCAction {
    if len(r.NextActions) == 0 {
        return nil
    }
    return &r.NextActions[0]
}

// RemoveCompletedAction removes the first action from the list
func (r *NPCResponse) RemoveCompletedAction() {
    if len(r.NextActions) > 0 {
        r.NextActions = r.NextActions[1:]
    }
}
