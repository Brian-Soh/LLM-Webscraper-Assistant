package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings" 
	"time"
)

const defaultBase = "http://localhost:11434"

type generateReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResp struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}
func Generate(model, prompt string) (string, error) {
	body, _ := json.Marshal(generateReq{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})

	req, err := http.NewRequest(http.MethodPost, defaultBase+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return "", errors.New("ollama error: " + string(b))
	}

	// Decode Ollama's returned JSON objects line by line and concatenate "response" fields.
	var full strings.Builder
	dec := json.NewDecoder(res.Body)
	for {
		var msg map[string]any
		if err := dec.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if s, ok := msg["response"].(string); ok {
			full.WriteString(s)
		} else if s, ok := msg["output"].(string); ok {
			full.WriteString(s)
		} else if s, ok := msg["message"].(string); ok {
			full.WriteString(s)
		}
	}

	text := strings.TrimSpace(full.String())
	if text == "" {
		return "", errors.New("ollama returned empty response")
	}
	return text, nil
}
