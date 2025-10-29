package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/ai-webscraper-go/ollama"
	"github.com/example/ai-webscraper-go/scraper"
)

// --- SCRAPE HANDLER ---

type ScrapeRequest struct {
	URL string `json:"url"`
}

type ScrapeResponse struct {
	URL          string `json:"url"`
	HTMLLength   int    `json:"htmlLength"`
	BodyLength   int    `json:"bodyLength"`
	Cleaned      string `json:"cleaned"`
	ErrorMessage string `json:"error,omitempty"`
}

func handleScrape(c *gin.Context) {
	var req ScrapeRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid 'url' in JSON body."})
		return
	}

	html, body, cleaned, err := scraper.ScrapeAndClean(req.URL)
	if err != nil {
		c.JSON(http.StatusBadGateway, ScrapeResponse{
			URL:          req.URL,
			HTMLLength:   len(html),
			BodyLength:   len(body),
			ErrorMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ScrapeResponse{
		URL:        req.URL,
		HTMLLength: len(html),
		BodyLength: len(body),
		Cleaned:    cleaned,
	})
}

// --- PARSE HANDLER ---

type ParseRequest struct {
	URL           string `json:"url,omitempty"`
	DOMContent    string `json:"domContent,omitempty"`
	Question      string `json:"question"`
	Model         string `json:"model,omitempty"`
	MaxChunkChars int    `json:"maxChunkChars,omitempty"` // default 100000
}

type ParseResponse struct {
	Model        string `json:"model"`
	Chunks       int    `json:"chunks"`
	Answer       string `json:"answer"`
	ErrorMessage string `json:"error,omitempty"`
}

func handleParse(c *gin.Context) {
	var req ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Question) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide 'question' and either 'domContent' or 'url'."})
		return
	}

	// Step 1: Retrieve and clean the page content
	content := strings.TrimSpace(req.DOMContent)
	if content == "" && strings.TrimSpace(req.URL) != "" {
		_, _, cleaned, err := scraper.ScrapeAndClean(req.URL)
		if err != nil {
			c.JSON(http.StatusBadGateway, ParseResponse{ErrorMessage: "Scrape failed: " + err.Error()})
			return
		}
		content = cleaned
	}

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No DOM content to parse. Provide 'domContent' or 'url'."})
		return
	}

	// Step 2: Chunk the content
	if req.MaxChunkChars <= 0 {
		req.MaxChunkChars = 100000
	}
	chunks := chunkString(content, req.MaxChunkChars)
	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = "gemma2:2b"
	}

	fmt.Println(">>> /api/parse started. Chunks:", len(chunks), "Model:", model)

	promptTmpl := "Extract information from the following content:\n\n%s\n\nQuestion: %s"
	var builder strings.Builder

	// Step 3: Send each chunk to Ollama
	for i, chunk := range chunks {
		fmt.Println(">>> Sending chunk", i+1, "of", len(chunks), "len:", len(chunk))

		done := make(chan struct{})
		var resp string
		var err error

		go func() {
			resp, err = ollama.Generate(model, fmt.Sprintf(promptTmpl, chunk, req.Question))
			close(done)
		}()

		select {
		case <-done:
			if err != nil {
				c.JSON(http.StatusBadGateway, ParseResponse{
					Model:        model,
					Chunks:       len(chunks),
					ErrorMessage: err.Error(),
				})
				return
			}
			builder.WriteString(resp)
			if i < len(chunks)-1 {
				builder.WriteString("\n")
			}
			fmt.Println(">>> Received chunk", i+1)
		case <-time.After(60 * time.Second):
			c.JSON(http.StatusGatewayTimeout, ParseResponse{
				Model:        model,
				Chunks:       len(chunks),
				ErrorMessage: "Ollama request timed out",
			})
			return
		}
	}

	c.JSON(http.StatusOK, ParseResponse{
		Model:  model,
		Chunks: len(chunks),
		Answer: builder.String(),
	})
}

// --- HELPERS ---

func chunkString(s string, max int) []string {
	if max <= 0 || len(s) <= max {
		return []string{s}
	}
	var out []string
	for i := 0; i < len(s); i += max {
		end := i + max
		if end > len(s) {
			end = len(s)
		}
		out = append(out, s[i:end])
	}
	return out
}
