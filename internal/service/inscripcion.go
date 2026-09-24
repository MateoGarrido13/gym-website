package service

import (
	"context"
	"fmt"

	db "PRACTICO_DOS/db/sqlc"
)

type InscripcionService struct {
	store InscripcionStore
}

func NewInscripcionService(store InscripcionStore) *InscripcionService {
	return &InscripcionService{store: store}
}

func (s *InscripcionService) Create(ctx context.Context, horarioID int32, in InscripcionCreateInput) (Inscripcion, error) {
	if horarioID == 0 || in.AlumnoID == 0 {
		return Inscripcion{}, fmt.Errorf("%w: horario_id y alumno_id son obligatorios", ErrInvalid)
	}

	cupo, err := s.store.GetCupoDeHorario(ctx, horarioID)
	if err != nil {
		return Inscripcion{}, mapNotFound(err)
	}
	if cupo.CantidadInscriptos >= cupo.MaxAlumnos {
		return Inscripcion{}, fmt.Errorf("%w: el horario alcanzó el máximo de alumnos", ErrConflict)
	}

	row, err := s.store.CreateInscripcion(ctx, db.CreateInscripcionParams{
		ClaseHorarioID: horarioID,
		AlumnoID:       in.AlumnoID,
	})
	if err != nil {
		return Inscripcion{}, err
	}
	return fromDBInscripcion(row), nil
}

func (s *InscripcionService) List(ctx context.Context, horarioID int32) ([]Inscripto, error) {
	if horarioID == 0 {
		return nil, fmt.Errorf("%w: horario_id es obligatorio", ErrInvalid)
	}

	rows, err := s.store.ListInscriptos(ctx, horarioID)
	if err != nil {
		return nil, err
	}
	out := make([]Inscripto, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBInscripto(row))
	}
	return out, nil
}

func (s *InscripcionService) Delete(ctx context.Context, horarioID, alumnoID int32) error {
	if horarioID == 0 || alumnoID == 0 {
		return fmt.Errorf("%w: horario_id y alumno_id son obligatorios", ErrInvalid)
	}
	return s.store.DeleteInscripcion(ctx, db.DeleteInscripcionParams{
		ClaseHorarioID: horarioID,
		AlumnoID:       alumnoID,
	})
}
