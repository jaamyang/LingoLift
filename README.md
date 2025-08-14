### LingoLift

A lightweight language-learning web app. Record or upload audio to get a transcript in the browser (Whisper via transformers.js). Click words to get context-aware explanations from an LLM and translate sentences via a Go backend.

#### Tech stack
- Frontend: Vue 3, Vite, Tailwind CSS, transformers.js (Whisper)
- Backend: Go (Gin)

#### Repository layout
```
lingolift_git/
├── backend/            # Go server (serves API + static build)
├── frontend/           # Vue app (Vite)
└── frontend_build.sh   # (legacy) build helper – paths may be outdated
```

### Prerequisites
- Node.js 18+ and npm
- Go 1.21+ (recommended)

Optional for LLM access:
- An LLM API key (if your chosen provider requires it)

### Quick start (local development)
Run backend API and the frontend dev server separately.

1) Backend (API)
```
go run ./backend
```
The server starts on http://localhost:8080 and exposes REST endpoints under `/api`.

2) Frontend (Vite dev)
```
cd frontend
npm install
npm run dev
```
Open the URL printed by Vite (typically http://localhost:5173). The backend has permissive CORS enabled, so the dev UI can call the API at http://localhost:8080.

### Production build and run
1) Build the frontend
```
cd frontend
npm install
npm run build
```
This produces `frontend/dist`.

2) Run the backend from the repository root so it can serve `frontend/dist` correctly:
```
go run ./backend
```
Then open http://localhost:8080. Static assets are served from `frontend/dist` and API remains under `/api`.

Note: The provided `backend/Dockerfile` and `frontend_build.sh` reference older paths (e.g., `backend-go/`). If you use Docker, adjust the paths in that file to the current layout or build/run locally as above.

### Configuration (environment variables)
The backend can call different LLM providers. Defaults target ModelScope Qwen.

- `LLM_BASE_URL` (optional): Base URL for chat/completions API. Default: `https://api-inference.modelscope.cn/v1/chat/completions`
- `LLM_MODEL` (optional): Model identifier. Default: `qwen/qwen3-30b-a3b`
- `LLM_API_KEY` (optional): API key for your LLM provider; if unset, Authorization header is omitted.

You can also override these per-request using the `llm_settings` object in request bodies (`baseUrl`, `model`, `apiKey`).

### API reference

Base URL (local): `http://localhost:8080`

- Health check: `GET /api/health`
```
HTTP/1.1 200 OK
{ "status": "healthy" }
```

- Explain a word in context: `POST /api/explain-word`
Request body:
```
{
  "word": "bark",
  "context": "I heard the dog bark loudly in the yard.",
  "audio_language": "English",
  "user_native_language": "English",
  "llm_settings": { "baseUrl": "...", "model": "...", "apiKey": "..." }
}
```
Response:
```
{
  "audio_language_explanation": "...",
  "native_language_explanation": "..."
}
```

- Translate a sentence: `POST /api/translate-sentence`
Request body:
```
{
  "sentence": "How are you today?",
  "target_language": "Spanish",
  "llm_settings": { "baseUrl": "...", "model": "...", "apiKey": "..." }
}
```
Response:
```
{ "translated_sentence": "¿Cómo estás hoy?" }
```

- Proxy a public asset (used by the frontend for models/wasm when needed): `GET /api/proxy?url=<ENCODED_URL>`
Response body is the proxied content with the original `Content-Type`.

### Frontend overview
Key components in `frontend/src/components/`:
- `AudioTranscriber.vue`: in-browser audio transcription using Whisper (transformers.js)
- `SettingsPage.vue`: configure model and preferences
- `TranscriptionHistory.vue`: view past transcriptions

Build tooling: Vite + Tailwind. Entry points are `frontend/src/main.js` and `frontend/src/App.vue`.

### Notes and troubleshooting
- If you see 404s for static assets when running the Go server, ensure you started it from the repository root after building the frontend so `frontend/dist` is discoverable.
- If your LLM provider requires an API key, set `LLM_API_KEY` or pass `llm_settings.apiKey`.
- Dockerfile paths may need updates to reflect `backend/` and `frontend/` directories.

### License
Unlicense license.