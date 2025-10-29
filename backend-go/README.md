# Go Backend

## Dev Setup

1. Install Go 1.22+ and Chrome/Chromium.
2. `go mod tidy`
3. `go run .` (server listens on `:8080`)

## Endpoints

### `POST /api/scrape`
```json
{ "url": "https://example.com" }
```
Response:
```json
{ "url":"...", "htmlLength":12345, "bodyLength":6789, "cleaned":"...text..." }
```

### `POST /api/parse`
```json
{ "url":"https://example.com", "question":"What are the store hours?" }
```
or
```json
{ "domContent":"...cleaned text...", "question":"Summarize the page" }
```
Response:
```json
{ "model":"llama3.2", "chunks": 1, "answer": "..." }
```

> Ensure **Ollama** is running locally.