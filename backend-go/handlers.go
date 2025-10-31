package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Brian-Soh/LLM-Webscraper-Assistant/backend-go/ollama"
	"github.com/Brian-Soh/LLM-Webscraper-Assistant/backend-go/scraper"
)

var promptTemplate string = "You are a helpful assistant. Answer the user's question using ONLY this chunk of page text.\n"+
			"If the chunk lacks the answer, say so briefly.\n\nQUESTION:\n%s\n\nCHUNK:\n%s"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provide a valid 'url'"})
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
	MaxChunkChars int    `json:"maxChunkChars,omitempty"`
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

	// Retrieve and clean the page content
	content := strings.TrimSpace(req.DOMContent)
	if content == "" && strings.TrimSpace(req.URL) != "" {
		_, _, cleaned, err := scraper.ScrapeAndClean(req.URL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Scrape failed: " + err.Error()})
			return
		}
		content = cleaned
	}

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No DOM content to parse. Provide omContent' or 'url'."})
		return
	}

	// Chunk the content
	if req.MaxChunkChars <= 0 {
		req.MaxChunkChars = 100000
	}
	chunks := chunkString(content, req.MaxChunkChars)

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = "gemma2:2b"
	}

	fmt.Println(">>> /api/parse started. Chunks:", len(chunks), "Model:", model)
	
	var builder strings.Builder
	// Send each chunk to Ollama
		for i, chunk := range chunks {
		prompt := fmt.Sprintf(
			promptTemplate,
			req.Question,
			chunk,
		)
		part, err := ollama.Generate(model, prompt)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "ollama: " + err.Error(), "chunksProcessed": i})
			return
		}
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(part)
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
