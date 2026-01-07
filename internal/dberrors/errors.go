package dberrors

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAlreadyExists   = errors.New("resource already exists")
	ErrNotFound        = errors.New("resource not found")
	ErrMissingRequired = errors.New("missing required field")
	ErrInvalidValue    = errors.New("invalid value")
	ErrInvalidRef      = errors.New("invalid reference")
)

func Handle(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	// Postgres-specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrAlreadyExists
		case "23502":
			return ErrMissingRequired
		case "23503":
			return ErrInvalidRef
		case "23514":
			return ErrInvalidValue
		}
	}
	return fmt.Errorf("database error: %w", err)
}
