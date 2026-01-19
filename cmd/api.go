package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/api/candles"
	"cmd/internal/api/sessions"
	"cmd/internal/api/web" // your htmx handlers

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config    config
	db        *pgx.Conn
	templates *template.Template
}

type config struct {
	addr string
	dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// ================================
	// Static files
	// ================================
	r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("web/js"))))
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))

	// ================================
	// Web/htmx routes (returns HTML)
	// ================================
	webHandler := web.NewHandler(repo.New(app.db), app.templates)

	// Pages
	r.Get("/", webHandler.Index)
	r.Get("/dashboard", webHandler.Dashboard)
	r.Get("/sessions-page", webHandler.SessionsPage)

	// htmx partials
	r.Route("/htmx", func(r chi.Router) {
		r.Get("/candles/table", webHandler.CandlesTable)
		r.Get("/candles/search", webHandler.SearchCandles)
		r.Delete("/candles/{id}", webHandler.DeleteCandle)
		r.Get("/sessions/list", webHandler.SessionsList)
		r.Get("/sessions/active", webHandler.ActiveSessions)
		r.Post("/sessions/{id}/toggle-favorite", webHandler.ToggleFavorite)
	})

	// ================================
	// JSON API (your existing handlers - unchanged)
	// ================================
	candleService := candles.NewService(repo.New(app.db))
	candleHandler := candles.NewHandler(candleService)

	r.Post("/candles", candleHandler.CreateCandle)
	r.Get("/candles/{id}", candleHandler.GetCandleByID)
	r.Get("/candles/params", candleHandler.GetCandle)
	r.Get("/candles/latest", candleHandler.GetLatestCandle)
	r.Get("/candles/recent", candleHandler.GetRecentCandles)
	// ... all your other candle routes stay the same ...

	sessionService := sessions.NewService(repo.New(app.db))
	sessionHandler := sessions.NewHandler(sessionService)

	r.Post("/sessions", sessionHandler.CreateSession)
	r.Get("/sessions/{id}", sessionHandler.GetSessionByID)
	// ... all your other session routes stay the same ...

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("Server has started at addr %v", app.config.addr)
	return srv.ListenAndServe()
}

func newApplication(templates *template.Template) *application {
	return &application{
		config: config{
			addr: ":8080",
			dbConfig: dbConfig{
				dsn: os.Getenv("DATABASE_URL"),
			},
		},
		templates: templates,
	}
}
