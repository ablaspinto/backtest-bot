package main

import (
	repo "cmd/internal/adapters/postgresql/sqlc"
	"cmd/internal/api/candles"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
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

	candleService := candles.NewService(repo.New(app.db))
	candleHandler := candles.NewHandler(candleService)
	r.Post("/candles", candleHandler.CreateCandle)
	r.Get("/candles/{id}", candleHandler.GetCandleByID)
	r.Get("/candles/params", candleHandler.GetCandle)
	r.Get("/candles/latest", candleHandler.GetLatestCandle)
	r.Get("/candles/recent", candleHandler.GetRecentCandles)
	r.Get("/candles/in-range", candleHandler.GetCandlesInRange)
	r.Get("/candles/timestamp/after", candleHandler.GetCandlesAfterTimestamp)
	r.Get("/candles/timestamp/before", candleHandler.GetCandlesBeforeTimestamp)
	r.Get("/candles/paginated", candleHandler.ListCandlesPaginated)
	r.Get("/candles/symbol", candleHandler.GetCandlesBySymbol)
	r.Get("/candles/count", candleHandler.CountCandles)
	r.Get("/candles/count/in-range", candleHandler.CountCandlesInRange)
	r.Get("/candles/stats", candleHandler.GetCandleStats)
	r.Get("/candles/stats/in-range", candleHandler.GetCandleStatsInRange)
	r.Get("/candles/volume", candleHandler.GetVolumeLeaders)
	r.Put("/candles", candleHandler.UpdateCandle)
	r.Delete("/candles/{id}", candleHandler.DeleteCandle)
	r.Delete("/candles/symbol", candleHandler.DeleteCandlesBySymbol)
	r.Delete("/candles/time-frame", candleHandler.DeleteCandlesByTimeFrame)
	r.Delete("/candles/old", candleHandler.DeleteOldCandles)
	r.Delete("/candles/in-range", candleHandler.DeleteCandlesInRange)
	r.Get("/candles/symbols/distinct", candleHandler.GetDistinctSymbols)
	r.Get("/candles/time-frame/distinct", candleHandler.GetDistinctTimeFrames)
	r.Get("/candles/pairs", candleHandler.GetSymbolTimeframePairs)
	r.Get("/candles/exists", candleHandler.CheckCandleExists)
	r.Get("/candles/oldest", candleHandler.GetOldestCandle)
	r.Get("/candles/gap", candleHandler.GetCandleGaps)

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

func newApplication() *application {
	return &application{
		config: config{
			addr: ":8080",
			dbConfig: dbConfig{
				dsn: os.Getenv("DATABASE_URL"),
			},
		},
	}
}
