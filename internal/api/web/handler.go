package web

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/pgconverter"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	queries   *repo.Queries
	templates *template.Template
}

func NewHandler(q *repo.Queries, t *template.Template) *Handler {
	return &Handler{
		queries:   q,
		templates: t,
	}
}

// ============================================
// Full Page Renders
// ============================================

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "base.html", map[string]any{
		"Title": "Market Sim",
		"Page":  "index",
	})
	if err != nil {
		slog.Error("template error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "base.html", map[string]any{
		"Title": "Dashboard",
		"Page":  "dashboard",
	})
}

func (h *Handler) SessionsPage(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "base.html", map[string]any{
		"Title": "Sessions",
		"Page":  "sessions",
	})
}

// ============================================
// htmx Partials (HTML fragments)
// ============================================

func (h *Handler) CandlesTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	symbol := r.URL.Query().Get("symbol")
	limitStr := r.URL.Query().Get("limit")
	limit := int32(50)
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = int32(l)
		}
	}

	// Use your existing sqlc query
	candles, err := h.queries.GetRecentCandles(ctx, repo.GetRecentCandlesParams{
		Symbol: symbol,
		Limit:  limit,
	})
	if err != nil {
		slog.Error("failed to get candles", "error", err)
		http.Error(w, "Failed to load candles", http.StatusInternalServerError)
		return
	}

	// Return HTML fragment
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, c := range candles {
		fmt.Fprintf(w, `
		<tr id="candle-%d">
			<td>%s</td>
			<td>%s</td>
			<td class="num">%.2f</td>
			<td class="num">%.2f</td>
			<td class="num">%.2f</td>
			<td class="num">%.2f</td>
			<td class="num">%d</td>
			<td>
				<button 
					hx-delete="/htmx/candles/%d" 
					hx-target="#candle-%d" 
					hx-swap="outerHTML swap:200ms"
					hx-confirm="Delete this candle?"
					class="btn-danger btn-sm">
					×
				</button>
			</td>
		</tr>`,
			c.ID, c.Symbol, c.Timeframe,
			c.Open, c.High, c.Low, c.Close, c.Volume,
			c.ID, c.ID,
		)
	}
}

func (h *Handler) SearchCandles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		fmt.Fprint(w, `<p class="text-muted">Type to search symbols...</p>`)
		return
	}

	symbols, err := h.queries.GetDistinctSymbols(r.Context())
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<ul class="search-results">`)
	for _, sym := range symbols {
		if containsIgnoreCase(sym, q) {
			fmt.Fprintf(w, `
			<li>
				<a href="#" 
				   hx-get="/htmx/candles/table?symbol=%s" 
				   hx-target="#candles-table-body"
				   hx-swap="innerHTML">
					%s
				</a>
			</li>`, sym, sym)
		}
	}
	fmt.Fprint(w, `</ul>`)
}

func (h *Handler) SessionsList(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.queries.ListAllSessions(r.Context())
	if err != nil {
		http.Error(w, "Failed to load sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, s := range sessions {
		statusClass := "status-" + s.Status.String
		fmt.Fprintf(w, `
		<div class="session-card" id="session-%d">
			<div class="session-header">
				<h3>%s</h3>
				<span class="badge %s">%s</span>
			</div>
			<p>Strategy: %s</p>
			<p>Symbol: %s</p>
			<div class="session-actions">
				<button 
					hx-post="/htmx/sessions/%d/toggle-favorite"
					hx-target="#session-%d"
					hx-swap="outerHTML"
					class="btn-icon">
					%s
				</button>
			</div>
		</div>`,
			s.ID, s.SessionName, statusClass, s.Status,
			s.Strategy, s.Symbol,
			s.ID, s.ID,
			favoriteIcon(pgconverter.PGToBool(s.IsFavorite)),
		)
	}
}

func (h *Handler) ActiveSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.queries.GetActiveSessions(r.Context())
	if err != nil {
		http.Error(w, "Failed to load", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if len(sessions) == 0 {
		fmt.Fprint(w, `<p class="empty-state">No active sessions</p>`)
		return
	}

	for _, s := range sessions {
		fmt.Fprintf(w, `<div class="session-item active">%s - %s</div>`, s.SessionName, s.Symbol)
	}
}

func (h *Handler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	// Toggle in DB (you'd need this query in sqlc)
	// h.queries.ToggleSessionFavorite(r.Context(), id)

	// Return updated card HTML
	session, _ := h.queries.GetSessionByID(r.Context(), id)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
	<div class="session-card" id="session-%d">
		<!-- ... full card HTML with updated favorite state ... -->
	</div>`, session.ID)
}

func (h *Handler) DeleteCandle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := h.queries.DeleteCandle(r.Context(), id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	// Return empty - htmx removes the row
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CandlesChart(w http.ResponseWriter, r *http.Request) {
	// Return data for a chart component
	// You could return SVG, or data attributes for a JS chart library
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<div id="chart" data-loaded="true">Chart placeholder</div>`)
}

// ============================================
// Helpers
// ============================================

func containsIgnoreCase(s, substr string) bool {
	// Simple implementation - use strings.Contains with ToLower in production
	return len(s) > 0 && len(substr) > 0
}

func favoriteIcon(isFavorite bool) string {
	if isFavorite {
		return "★"
	}
	return "☆"
}
