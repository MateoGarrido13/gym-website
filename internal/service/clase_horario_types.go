package service

import (
	"fmt"
	"strings"
	"time"

	db "PRACTICO_DOS/db/sqlc"
)

type ClaseHorario struct {
	ID         int32  `json:"id"`
	ClaseID    int32  `json:"clase_id"`
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
}

type ClaseHorarioCreateInput struct {
	DiaSemana  string `json:"dia_semana"`
	HoraInicio string `json:"hora_inicio"`
}

func fromDBClaseHorario(row db.ClaseHorario) ClaseHorario {
	return ClaseHorario{
		ID:         row.ID,
		ClaseID:    row.ClaseID,
		DiaSemana:  row.DiaSemana,
		HoraInicio: row.HoraInicio.Format("15:04:05"),
	}
}

func parseHoraInicio(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, layout := range []string{"15:04:05", "15:04"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("%w: hora_inicio debe ser HH:MM o HH:MM:SS", ErrInvalid)
}

func normalizeDiaSemana(s string) (string, error) {
	dia := strings.ToLower(strings.TrimSpace(s))
	switch dia {
	case "lunes", "martes", "miercoles", "jueves", "viernes", "sabado", "domingo":
		return dia, nil
	default:
		return "", fmt.Errorf("%w: dia_semana debe ser lunes, martes, miercoles, jueves, viernes, sabado o domingo", ErrInvalid)
	}
}
