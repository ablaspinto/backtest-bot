# Backtesting Automated









## Work to be done

### Candles API
  - setup mock data / load mock data

  - Test api endpoints

### Sessions API
  - setup mock data / load mock data

  - Test api endpoints

### User Configs API
  - setup mock data / load mock data

  - Test api endpoints

### Front end


  - Fetch Backend candle data, load 1d data first






# htmx + Go Architecture

## Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                          BROWSER                                │
│                    http://localhost:8080                        │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                       GO SERVER (:8080)                         │
│                                                                 │
│  cmd/main.go    → Starts server, loads templates, connects DB   │
│  cmd/api.go     → Defines all routes                            │
└─────────────────────────────────┬───────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                    POSTGRESQL (Docker :5432)                    │
│                                                                 │
│  Tables: candles, sessions                                      │
└─────────────────────────────────────────────────────────────────┘
```

---

## Route Types

| Route Type | Example | Handler | Returns |
|------------|---------|---------|---------|
| Static Files | `/css/*`, `/js/*` | FileServer | CSS, JS, images |
| HTML Pages | `/`, `/dashboard` | web.Handler | Full HTML page |
| htmx Partials | `/htmx/candles/table` | web.Handler | HTML fragment |
| JSON API | `/candles/recent` | candles.handler | JSON |

---

## Request Flow

### 1. Initial Page Load

```
Browser                         Server
   │                              │
   │  GET /dashboard              │
   │─────────────────────────────▶│
   │                              │
   │                              │  Renders base.html
   │                              │  + embeds dashboard.html
   │                              │
   │◀───── Full HTML Page ────────│
   │                              │
   │  Browser displays page       │
   │  htmx.js loaded              │
   │                              │
```

### 2. htmx Dynamic Loading

```
Browser                         Server                      Database
   │                              │                            │
   │  htmx sees:                  │                            │
   │  hx-get="/htmx/candles/table"│                            │
   │  hx-trigger="load"           │                            │
   │                              │                            │
   │  GET /htmx/candles/table     │                            │
   │─────────────────────────────▶│                            │
   │                              │  SELECT * FROM candles     │
   │                              │───────────────────────────▶│
   │                              │◀────── rows ───────────────│
   │                              │                            │
   │◀──── <tr>...</tr> ───────────│                            │
   │      (HTML fragment)         │                            │
   │                              │                            │
   │  htmx swaps into DOM         │                            │
   │  (#candles-table-body)       │                            │
```

### 3. User Interaction (Search)

```
Browser                         Server                      Database
   │                              │                            │
   │  User types "ES" in search   │                            │
   │                              │                            │
   │  htmx sees:                  │                            │
   │  hx-get="/htmx/candles/search"                            │
   │  hx-trigger="input changed delay:300ms"                   │
   │                              │                            │
   │  GET /htmx/candles/search?q=ES                            │
   │─────────────────────────────▶│                            │
   │                              │  SELECT DISTINCT symbol    │
   │                              │  WHERE symbol ILIKE '%ES%' │
   │                              │───────────────────────────▶│
   │                              │◀────── symbols ────────────│
   │                              │                            │
   │◀──── <ul><li>ES</li></ul> ───│                            │
   │                              │                            │
   │  htmx swaps into             │                            │
   │  #search-results             │                            │
```

---

## File Structure

```
~/backtest-bot/
│
├── cmd/
│   ├── main.go                      # Entry point
│   ├── api.go                       # Routes + middleware
│   │
│   └── internal/
│       ├── adapters/postgresql/
│       │   ├── sqlc/                # Generated by sqlc
│       │   │   ├── queries.sql.go
│       │   │   └── models.go
│       │   └── migrations/          # SQL migrations
│       │       ├── 001_candles.sql
│       │       └── 002_sessions.sql
│       │
│       └── api/
│           ├── candles/
│           │   ├── handler.go       # JSON responses
│           │   ├── service.go
│           │   └── params.go
│           │
│           ├── sessions/
│           │   ├── handler.go       # JSON responses
│           │   ├── service.go
│           │   └── params.go
│           │
│           └── web/
│               └── handler.go       # HTML responses (htmx)
│
├── web/
│   ├── assets/                      # Images, fonts
│   ├── css/
│   │   └── style.css
│   ├── js/                          # Optional custom JS
│   │
│   └── templates/
│       ├── layouts/
│       │   └── base.html            # Main layout wrapper
│       │
│       └── pages/
│           ├── index.html           # {{define "index"}}
│           ├── dashboard.html       # {{define "dashboard"}}
│           └── sessions.html        # {{define "sessions"}}
│
├── docker-compose.yaml              # PostgreSQL container
├── Makefile                         # Dev commands
├── .air.toml                        # Hot reload config
├── go.mod
└── go.sum
```

---

## Template System

### base.html (Layout)

```html
{{define "base.html"}}
<!DOCTYPE html>
<html>
<head>
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <link rel="stylesheet" href="/css/style.css">
</head>
<body>
    <nav>...</nav>
    
    <main>
        {{if eq .Page "index"}}
            {{template "index" .}}
        {{else if eq .Page "dashboard"}}
            {{template "dashboard" .}}
        {{else if eq .Page "sessions"}}
            {{template "sessions" .}}
        {{end}}
    </main>
</body>
</html>
{{end}}
```

### Page Templates

```html
<!-- index.html -->
{{define "index"}}
<h1>Home</h1>
<div hx-get="/htmx/candles/table" hx-trigger="load">
    Loading...
</div>
{{end}}

<!-- dashboard.html -->
{{define "dashboard"}}
<h1>Dashboard</h1>
<table>
    <tbody hx-get="/htmx/candles/table" hx-trigger="load">
        Loading...
    </tbody>
</table>
{{end}}
```

---

## Handler Comparison

### JSON Handler (existing)

```go
// cmd/internal/api/candles/handler.go
func (h *handler) GetRecentCandles(w http.ResponseWriter, r *http.Request) {
    candles, _ := h.service.GetRecentCandles(ctx, params)
    json.WriteSuccess(w, http.StatusOK, candles)  // Returns JSON
}
```

**Response:**
```json
[{"id":1,"symbol":"ES","open":4500.25,"high":4510.00}]
```

### htmx Handler (new)

```go
// cmd/internal/api/web/handler.go
func (h *Handler) CandlesTable(w http.ResponseWriter, r *http.Request) {
    candles, _ := h.queries.GetRecentCandles(ctx, params)
    
    w.Header().Set("Content-Type", "text/html")
    for _, c := range candles {
        fmt.Fprintf(w, `<tr><td>%s</td><td>%.2f</td></tr>`, 
            c.Symbol, c.Open)  // Returns HTML
    }
}
```

**Response:**
```html
<tr><td>ES</td><td>4500.25</td></tr>
```

---

## htmx vs Traditional SPA

### Traditional SPA (React/Vue)

```
┌────────────┐         ┌────────────┐         ┌────────────┐
│  Browser   │  JSON   │   API      │         │  Database  │
│  (React)   │◀───────▶│  Server    │◀───────▶│            │
└────────────┘         └────────────┘         └────────────┘
     │
     │  JavaScript converts
     │  JSON → HTML
     │
     ▼
┌────────────┐
│  npm       │
│  webpack   │
│  node_modules
└────────────┘
```

### htmx Way

```
┌────────────┐         ┌────────────┐         ┌────────────┐
│  Browser   │  HTML   │   Go       │         │  Database  │
│  (htmx)    │◀───────▶│  Server    │◀───────▶│            │
└────────────┘         └────────────┘         └────────────┘
                            │
                            │  Server renders HTML
                            │  (no build step)
                            │
                       ┌────────────┐
                       │  Just Go   │
                       │  go run .  │
                       └────────────┘
```

---

## Common htmx Attributes

| Attribute | Purpose | Example |
|-----------|---------|---------|
| `hx-get` | GET request | `hx-get="/htmx/data"` |
| `hx-post` | POST request | `hx-post="/htmx/create"` |
| `hx-delete` | DELETE request | `hx-delete="/htmx/item/1"` |
| `hx-trigger` | When to fire | `hx-trigger="load"`, `"click"`, `"input delay:300ms"` |
| `hx-target` | Where to put response | `hx-target="#results"` |
| `hx-swap` | How to swap | `hx-swap="innerHTML"`, `"outerHTML"`, `"beforeend"` |
| `hx-confirm` | Confirmation dialog | `hx-confirm="Delete?"` |
| `hx-indicator` | Loading indicator | `hx-indicator="#spinner"` |

---

## Makefile Commands

```bash
make dev            # Start docker + air (hot reload)
make run            # Run the server
make up             # Start docker containers
make down           # Stop docker containers
make migrate-up     # Run all migrations
make migrate-down   # Rollback last migration
make db-shell       # Open psql shell
```

---

## Quick Start

```bash
# 1. Start database
make up

# 2. Run migrations
make migrate-up

# 3. Start server with hot reload
make dev

# 4. Open browser
open http://localhost:8080
```
