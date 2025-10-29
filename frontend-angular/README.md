# Angular Frontend
## Dev Setup

```bash
npm i -g @angular/cli
npm i
npm run start
```

This uses `proxy.conf.json` to forward `/api` to `http://localhost:8080`.

## Workflow
- Scrapes a URL via `POST /api/scrape`, displays "cleaned" text.
- Sends questions + content to `POST /api/parse` and displays the LLM response.