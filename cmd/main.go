package main

import (
	"cmd/internal/env"
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	//loader := data.NewLoader()
	//bars, err := loader.LoadSingleFile("CME_ESH2001.csv")
	//smaObj := indicators.SMA(bars, 10)
	//fmt.Printf("SMA BARS: %v\n", smaObj)
	//emaObj := indicators.EMA(bars, 10)
	//fmt.Printf("EMA BARS: %v\n", emaObj)
	//rsiObj := indicators.RSI(bars, 14)
	//fmt.Printf("RSI : %v\n", rsiObj)
	//macd := indicators.MACD(bars)
	//fmt.Printf("MACD HISTOGRAM: %v\n", macd)
	ctx := context.Background()
	cfg := config{
		addr: ":8080",
		dbConfig: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=marketsim sslmode=disable"),
		},
	}
	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// database
	conn, err := pgx.Connect(ctx, cfg.dbConfig.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("Connected to database", "dsn", cfg.dbConfig.dsn)
	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		logger.Info("server has failed to start %v\n", err)
		os.Exit(1)
	}
}
