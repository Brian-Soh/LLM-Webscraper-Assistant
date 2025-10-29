# LLM Webscraper Assistant

In this project I recreated my AI Webscraper project from Python to a full-stack Go (Gin) + Angular/TypeScript application with the goal of learning a new tech stack and reviewing my full-stack development skills.

The system allows users to scrape live websites, clean and normalize their content, and query a Large Language Model (LLM) through Ollama’s REST API to generate context-aware responses; all from a responsive, browser-based interface.

The backend, written in Go, uses Gin for routing, Chromedp for headless web scraping, and GoQuery for DOM parsing. The frontend, built in Angular/TypeScript, provides a reactive user interface with live query execution, model selection, and real-time updates. CORS and an Angular proxy configuration ensure seamless and secure communication between both layers.

## Quickstart
1. **Ollama**: Install & run locally, and pull a model: `ollama pull llama3.2`
2. **Backend**:
   ```bash
   cd backend-go
   go mod tidy
   go run .
   ```
   Server listens on `http://localhost:8080`

3. **Frontend**:
   ```bash
   cd ../frontend-angular
   npm i -g @angular/cli
   npm i
   npm start
   ```
   This opens `http://localhost:4200`, proxied to the backend.