package ollama

import (
	"bufio"
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

type ndjsonResp struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func Generate(model, prompt string) (response string, err error) {
	body, _ := json.Marshal(generateReq{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	})

	req, err := http.NewRequest(http.MethodPost, defaultBase+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 2 * time.Minute}
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

	// Use a scanner to allow for response streaming
	scanner := bufio.NewScanner(res.Body)

	// runs perpetually until EOF is received
	for scanner.Scan() {
		var piece ndjsonResp
		err = json.Unmarshal(scanner.Bytes(), &piece)
		if err == nil {
			full.WriteString(piece.Response)
		}
	}
	response = full.String()

	if response == "" {
		return "", errors.New("ollama returned empty response")
	}

	return response, nil

}
