package main

import (
	"cmd/internal/env"
	"context"
	"html/template"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Templates
	templates, err := template.ParseGlob("web/templates/layouts/*.html")
	if err != nil {
		logger.Error("failed to parse layout templates", "error", err)
		os.Exit(1)
	}
	templates, err = templates.ParseGlob("web/templates/pages/*.html")
	if err != nil {
		logger.Error("failed to parse page templates", "error", err)
		os.Exit(1)
	}

	// Database
	dsn := env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=marketsim sslmode=disable")
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	logger.Info("Connected to database")

	app := &application{
		config: config{
			addr:     ":8080",
			dbConfig: dbConfig{dsn: dsn},
		},
		db:        conn,
		templates: templates,
	}

	if err := app.run(app.mount()); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
