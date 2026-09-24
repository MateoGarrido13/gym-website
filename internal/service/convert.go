package service

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalid  = errors.New("datos inválidos")
	ErrNotFound = errors.New("no encontrado")
	ErrConflict = errors.New("conflicto")
)

func fromNullString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func fromNullInt32(v sql.NullInt32) *int32 {
	if !v.Valid {
		return nil
	}
	n := v.Int32
	return &n
}

func fromNullTime(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	s := v.Time.UTC().Format(time.RFC3339)
	return &s
}

func fromNullTimeDate(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	s := v.Time.UTC().Format(time.DateOnly)
	return &s
}

func toNullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func toNullInt32(v *int32) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *v, Valid: true}
}

func parseDate(s string) (time.Time, error) {
	return time.Parse(time.DateOnly, strings.TrimSpace(s))
}

func toNullDate(v *string) (sql.NullTime, error) {
	if v == nil || strings.TrimSpace(*v) == "" {
		return sql.NullTime{}, nil
	}
	t, err := parseDate(*v)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}
