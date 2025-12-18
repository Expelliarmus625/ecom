package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/expelliarmus625/ecom/internal/env"
	"github.com/jackc/pgx/v5"
)


func main() {
	ctx := context.Background()

	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=ecom sslmode=disable"),
		},
	}

	//Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	api := application{
		config: cfg,
		db: conn,
	}

	logger.Info("Connected to Database", "dsn", cfg.db.dsn)

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}
}
