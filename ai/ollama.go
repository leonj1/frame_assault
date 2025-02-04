package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/http"
    "time"
)

// Default timeout for HTTP requests
const defaultTimeout = 30 * time.Second

// OllamaClient handles communication with the Ollama API
type OllamaClient struct {
    host    string
    model   string
    timeout time.Duration
}

// OllamaRequest represents the request body for Ollama API
type OllamaRequest struct {
    Model     string         `json:"model"`
    Prompt    string         `json:"prompt"`
    Stream    bool           `json:"stream"`
    MaxTokens int           `json:"max_tokens,omitempty"`
    Context   *GameContext  `json:"context,omitempty"`
}

// OllamaResponse represents the response from Ollama API
type OllamaResponse struct {
    Model    string `json:"model"`
    Response string `json:"response"`
    Done     bool   `json:"done"`
    Error    string `json:"error,omitempty"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(host, model string) *OllamaClient {
    return &OllamaClient{
        host:    host,
        model:   model,
        timeout: defaultTimeout,
    }
}

// SetTimeout sets a custom timeout for HTTP requests
func (c *OllamaClient) SetTimeout(timeout time.Duration) {
    c.timeout = timeout
}

// GenerateResponse sends a prompt to Ollama and returns the response
func (c *OllamaClient) GenerateResponse(prompt string) (string, error) {
    return c.GenerateResponseWithContext(prompt, nil)
}

// GenerateResponseWithContext sends a prompt with game context to Ollama
func (c *OllamaClient) GenerateResponseWithContext(prompt string, context *GameContext) (string, error) {
    // Prepare request body
    reqBody := OllamaRequest{
        Model:   c.model,
        Prompt:  prompt,
        Stream:  false,
        Context: context,
    }
    
    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %v", err)
    }
    
    // Create HTTP request
    url := fmt.Sprintf("http://%s/api/generate", c.host)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    
    // Send request with timeout
    client := &http.Client{
        Timeout: c.timeout,
    }
    resp, err := client.Do(req)
    if err != nil {
        if err, ok := err.(net.Error); ok && err.Timeout() {
            return "", fmt.Errorf("request timed out after %v: %v", c.timeout, err)
        }
        return "", fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()
    
    // Check HTTP status code
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
    }
    
    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response: %v", err)
    }
    
    // Parse response
    var ollamaResp OllamaResponse
    if err := json.Unmarshal(body, &ollamaResp); err != nil {
        return "", fmt.Errorf("error parsing response: %v", err)
    }
    
    // Check for API error
    if ollamaResp.Error != "" {
        return "", fmt.Errorf("ollama API error: %s", ollamaResp.Error)
    }
    
    return ollamaResp.Response, nil
}

// GetStrategicAdvice generates strategic advice based on the current game context
func (c *OllamaClient) GetStrategicAdvice(context *GameContext) (string, error) {
    prompt := context.FormatPrompt()
    return c.GenerateResponseWithContext(prompt, context)
}

// GetNPCResponse generates and parses an NPC's next actions and state
func (c *OllamaClient) GetNPCResponse(context *GameContext, npc *ComputerUser) (*NPCResponse, error) {
    prompt := FormatNPCPrompt(context, npc)
    response, err := c.GenerateResponseWithContext(prompt, context)
    if err != nil {
        return nil, fmt.Errorf("failed to generate response: %v", err)
    }
    
    npcResponse, err := ParseOllamaResponse(response)
    if err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }
    
    if err := npcResponse.ValidateResponse(); err != nil {
        return nil, fmt.Errorf("invalid response: %v", err)
    }
    
    return npcResponse, nil
}
