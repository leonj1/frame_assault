package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// OllamaClient handles communication with the Ollama API
type OllamaClient struct {
    host  string
    model string
}

// OllamaRequest represents the request body for Ollama API
type OllamaRequest struct {
    Model    string `json:"model"`
    Prompt   string `json:"prompt"`
    Stream   bool   `json:"stream"`
    MaxTokens int   `json:"max_tokens,omitempty"`
}

// OllamaResponse represents the response from Ollama API
type OllamaResponse struct {
    Model     string `json:"model"`
    Response  string `json:"response"`
    Done      bool   `json:"done"`
    Error     string `json:"error,omitempty"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(host, model string) *OllamaClient {
    return &OllamaClient{
        host:  host,
        model: model,
    }
}

// GenerateResponse sends a prompt to Ollama and returns the response
func (c *OllamaClient) GenerateResponse(prompt string) (string, error) {
    // Prepare request body
    reqBody := OllamaRequest{
        Model:  c.model,
        Prompt: prompt,
        Stream: false,
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
    
    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()
    
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
