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
    Model     string `json:"model"`
    Prompt    string `json:"prompt"`
    Stream    bool   `json:"stream"`
    MaxTokens int    `json:"max_tokens,omitempty"`
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
