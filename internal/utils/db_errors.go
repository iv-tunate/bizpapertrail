package utils

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5"
)

func ParseDbError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, http.StatusText(200)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusNotFound, "Requested resource not found"
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return http.StatusNotFound, "The requested resource or dependency does not exist."
		case "23505":
			return http.StatusConflict, "This record or relationship already exists."
		case "23514":
			return http.StatusBadRequest, "The provided data violates system validation rules."
		case "22001":
			return http.StatusBadRequest, "One or more text fields exceed the allowed character limit."
		}
	}

	return http.StatusInternalServerError, "An unexpected database error occurred."
}