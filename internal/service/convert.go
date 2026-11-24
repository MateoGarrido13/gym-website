package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	db "PRACTICO_DOS/db/sqlc"
)

var (
	ErrInvalid  = errors.New("datos inválidos")
	ErrNotFound = errors.New("no encontrado")
	ErrConflict = errors.New("conflicto")

	// ErrUnauthorized cubre cualquier fallo de login o sesión inválida. El
	// mensaje es genérico a propósito: no dice si falló el email, la
	// contraseña o el rol.
	ErrUnauthorized = errors.New("credenciales inválidas")
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

// parseTipoPlan valida que el string sea uno de los valores permitidos del enum tipo_plan.
func ParseTipoPlan(s string) (db.EnumTipoPlan, error) {
	s = strings.TrimSpace(s)
	switch s {
	case "Mensual":
		return db.EnumTipoPlanMensual, nil
	case "Trimestral":
		return db.EnumTipoPlanTrimestral, nil
	case "Anual":
		return db.EnumTipoPlanAnual, nil
	default:
		return "", fmt.Errorf("tipo_plan invalido: %q (valores permitidos: Mensual, Trimestral, Anual)", s)
	}
}

// fromNullEnumTipoPlan convierte un NullEnumTipoPlan de la base de datos a *db.EnumTipoPlan del servicio.
func fromNullEnumTipoPlan(v db.NullEnumTipoPlan) *db.EnumTipoPlan {
	if !v.Valid {
		return nil
	}
	t := v.EnumTipoPlan
	return &t
}

// toNullEnumTipoPlan convierte un *db.EnumTipoPlan del servicio a NullEnumTipoPlan para la base de datos.
func toNullEnumTipoPlan(v *db.EnumTipoPlan) db.NullEnumTipoPlan {
	if v == nil {
		return db.NullEnumTipoPlan{}
	}
	return db.NullEnumTipoPlan{EnumTipoPlan: *v, Valid: true}
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if t, err := time.Parse(time.DateOnly, s); err == nil {
		return t, nil
	}
	return time.Parse("02-01-2006", s)
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
