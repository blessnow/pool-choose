# AGENTS.md

## Purpose
This file helps AI coding agents understand the `yuchi` repo quickly and work productively without guessing the architecture or dev commands.

## Repo overview
- `backend/`: Go Gin backend service.
- `frontend/`: React + TypeScript + Vite frontend.
- `docker-compose.yml`: local full-stack dev environment.
- `data/`: SQLite database and persisted data; `docker-compose` mounts this directory into the backend container.

## Key concepts
- Backend serves API under `/api` and also serves the frontend static files at `/`.
- Frontend uses `frontend/src/services/api.ts` as the main API client.
- Auth is password-based and uses an `Authorization` header token on requests.

## Run commands
- Backend dev: `cd backend && go mod tidy && go run cmd/server/main.go`
- Frontend dev: `cd frontend && npm install && npm run dev`
- Full stack: `docker-compose up -d`

## Important files
- `backend/cmd/server/main.go`: backend entrypoint, routes, CORS, static file service, env var defaults.
- `backend/internal/handlers/`: HTTP handlers, auth middleware, business endpoints.
- `backend/internal/repository/`: database initialization and persistence.
- `backend/internal/models/`: data models.
- `frontend/src/services/api.ts`: REST client and API contract.
- `frontend/src/components/`: UI components.
- `frontend/src/hooks/useStocks.ts`: stock data hook.

## API conventions
- Base backend API path: `http://localhost:8080/api` by default.
- Key API routes:
  - `POST /api/auth/login`
  - `POST /api/auth/logout`
  - `GET /api/stocks`
  - `GET /api/stocks/:code`
  - `POST /api/stocks`
  - `PUT /api/stocks/:code`
  - `DELETE /api/stocks/:code`
  - `GET /api/quotes?codes=`
  - `GET /api/company-summaries?codes=`
  - `GET /api/chart/:code?period=`
  - `GET /api/cycle-insight`
  - `PUT /api/cycle-insight`
  - `GET /api/positions`
  - `POST /api/positions`
  - `PUT /api/positions/:id`
  - `DELETE /api/positions/:id`

## Environment variables
- `PASSWORD`: login password, default `dayuchi` if not set in `docker-compose.yml`.
- `DB_PATH`: SQLite path, default `./data/cycle_stock.db`.
- `PORT`: backend port, default `8080`.
- `VITE_API_BASE`: frontend API base URL when running dev mode.

## Agent guidance
- Prefer working with the existing split between frontend and backend.
- If modifying API shapes, update both `frontend/src/services/api.ts` and the backend handler signatures.
- Do not assume a test suite exists unless added; focus on manual `npm run lint` and Go build/run validation.
- Preserve the SQLite persistence model and the mounted `data/` directory in `docker-compose.yml`.
- The repository is a personal investment research system; avoid changing text into financial advice.

## Notes
- The backend uses Sina Finance APIs for real-time quotes and macro data from public sources.
- Frontend is built with Vite + Tailwind; the `frontend/README.md` is generic scaffold documentation and not repo-specific.
