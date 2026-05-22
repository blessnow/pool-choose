# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

周期股投资管理系统 (Cyclical Stock Investment Management System) - A full-stack web application for managing cyclical stock investments with real-time quotes, K-line charts, and position tracking.

## Tech Stack

- **Backend**: Go 1.21+, Gin, GORM, SQLite
- **Frontend**: React 19, TypeScript, Vite, TailwindCSS, ECharts
- **Data Sources**: Tencent Finance API (real-time quotes), Sina Finance API (K-line data), Eastmoney API (valuation/PE/PB data)
- **Deployment**: Docker containerized deployment

## Development Commands

### Backend (from `backend/` directory)
```bash
go mod tidy
go run cmd/server/main.go
```
Backend runs on port 8080 by default.

### Frontend (from `frontend/` directory)
```bash
npm install
npm run dev      # Development server on port 5173
npm run build    # Production build
npm run lint     # ESLint validation
```

### Full Stack
```bash
docker-compose up -d
```
Access at http://localhost (frontend on port 80, backend on port 8080).

## Architecture

### Backend Structure
- `cmd/server/main.go` - Entry point, route definitions, CORS config, static file serving
- `internal/handlers/` - HTTP handlers (auth.go, stocks.go, cycle.go)
- `internal/services/` - External API clients and business logic: `quote_service.go` (Tencent quotes + Sina K-line), `valuation_band_service.go` (PE/PB historical bands), `macro_service.go` (CPI/PPI from Eastmoney), `stock_search_service.go` (code/name lookup)
- `internal/models/` - GORM data models (models.go)
- `internal/repository/` - Database access layer: `db.go` (init), `seed.go` (initial stock pool seeding on first run)

### Frontend Structure
- `src/services/api.ts` - REST client, all API type definitions and calls
- `src/hooks/useStocks.ts` - Stock data hook with auto-refresh
- `src/components/` - UI components (StockCard, StockDetail, FilterBar, CycleSection, etc.)
- `src/App.tsx` - Main application component

### Key Architectural Notes

1. **Backend serves frontend**: In production, the Go backend serves static files at `/` and `/assets/*`, with API routes under `/api`.

2. **Authentication**: Password-based auth using `Authorization` header token. Default password: `dayuchi`.

3. **Data persistence**: SQLite database stored in `data/` directory, mounted in Docker.

4. **External APIs**: Real-time quotes from Tencent Finance (`web.sqt.gtimg.cn`), K-line data from Sina Finance, valuation data from Eastmoney. These APIs have rate limits - implement appropriate caching. Monthly K-line data is aggregated from daily data in `quote_service.go:aggregateToMonthly()`.

## Environment Variables

- `PASSWORD` - Login password (default: `dayuchi`)
- `DB_PATH` - SQLite database path (default: `./data/cycle_stock.db`)
- `PORT` - Backend port (default: `8080`)
- `VITE_API_BASE` - Frontend API base URL for dev mode

## API Contract

When modifying API shapes, update both `frontend/src/services/api.ts` and corresponding backend handlers.

All API responses use `{ ok: bool, data?: T, error?: string }`. Auth token is passed as `Authorization` header (plain token, not Bearer).

Key routes:
- `POST /api/auth/login`, `POST /api/auth/logout`
- `GET/POST/PUT/DELETE /api/stocks`, `/api/stocks/:code`
- `GET /api/stock-search?q=`, `GET /api/stock-info?code=` (lookup helpers backed by `stock_search_service.go`)
- `GET /api/quotes?codes=`, `GET /api/company-summaries?codes=`, `GET /api/chart/:code?period=daily|weekly|monthly`, `GET /api/valuation/:code`
- `GET /api/valuation-band?code=&metric=pe_ttm|pb&years=1-10`
- `GET/PUT /api/cycle-insight`
- `GET/POST/PUT/DELETE /api/positions`, `/api/positions/:id`; `GET /api/positions?code=` filters by stock

## Key Implementation Details

- **`SinaQuoteClient` naming**: Despite the name, `GetQuotes` uses the Tencent Finance API (`web.sqt.gtimg.cn`), not Sina. Sina is only used for K-line data.
- **`CycleInsight` JSON serialization**: `macroCards`, `bars`, and `sources` are slice fields with `gorm:"-"` — they are stored in separate `*JSON` string columns and manually marshaled/unmarshaled in `handlers/cycle.go`.
- **`valuation_band_service.go`**: Computes historical PE TTM and PB bands from K-line + quarterly financial data (EPS/BPS from Eastmoney). Adds percentile tracks (p10/p30/p50/p70/p90) across the full date range. A股 financials are cumulative YTD, not quarterly — `getTTMEPS` sums the last 4 report values as an approximation.
- **`UpdatePosition` / `DeletePosition`**: Currently stub implementations — they ignore the `:id` param. Functional position management is read/create only.

## Notes

- This is a personal investment research system - avoid changing text into financial advice.
- No test suite exists - validate with `npm run lint` (frontend) and `go build` (backend).
- The `frontend/README.md` is generic Vite scaffold documentation, not project-specific.
- Stock code format: 6-digit codes (e.g., `000707`, `600000`). Codes starting with `6` or `5` are Shanghai (sh prefix), others are Shenzhen (sz prefix).
