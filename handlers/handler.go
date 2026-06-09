package handlers

import (
	"log/slog"

	"github.com/iv-tunate/bizpapertrail/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *database.Queries
	Pool *pgxpool.Pool
	Logger *slog.Logger
}

func NEW(db *database.Queries, pool *pgxpool.Pool, logger *slog.Logger) *Handler{
	return  &Handler{DB: db, Pool: pool, Logger: logger}
}