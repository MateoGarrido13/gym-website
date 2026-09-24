package service

import (
	"context"
	"fmt"

	db "PRACTICO_DOS/db/sqlc"
)

type RutinaEjercicioService struct {
	store RutinaEjercicioStore
}

func NewRutinaEjercicioService(store RutinaEjercicioStore) *RutinaEjercicioService {
	return &RutinaEjercicioService{store: store}
}

func (s *RutinaEjercicioService) Create(ctx context.Context, rutinaID int32, in RutinaEjercicioCreateInput) (RutinaEjercicio, error) {
	if rutinaID == 0 || in.EjercicioID == 0 {
		return RutinaEjercicio{}, fmt.Errorf("%w: rutina_id y ejercicio_id son obligatorios", ErrInvalid)
	}
	if err := validRutinaEjercicioNumeros(in.Orden, in.Series, in.Repeticiones, in.DescansoSegundos); err != nil {
		return RutinaEjercicio{}, err
	}

	err := s.store.AddEjercicioARutina(ctx, db.AddEjercicioARutinaParams{
		RutinaID:         rutinaID,
		EjercicioID:      in.EjercicioID,
		Orden:            in.Orden,
		Series:           toNullInt32(in.Series),
		Repeticiones:     toNullInt32(in.Repeticiones),
		DescansoSegundos: toNullInt32(in.DescansoSegundos),
	})
	if err != nil {
		return RutinaEjercicio{}, err
	}

	return RutinaEjercicio{
		RutinaID:         rutinaID,
		EjercicioID:      in.EjercicioID,
		Orden:            in.Orden,
		Series:           in.Series,
		Repeticiones:     in.Repeticiones,
		DescansoSegundos: in.DescansoSegundos,
	}, nil
}

func (s *RutinaEjercicioService) List(ctx context.Context, rutinaID int32) ([]RutinaEjercicio, error) {
	if rutinaID == 0 {
		return nil, fmt.Errorf("%w: rutina_id es obligatorio", ErrInvalid)
	}

	rows, err := s.store.ListRutinaEjercicios(ctx, rutinaID)
	if err != nil {
		return nil, err
	}
	out := make([]RutinaEjercicio, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBRutinaEjercicio(row))
	}
	return out, nil
}

func (s *RutinaEjercicioService) Update(ctx context.Context, rutinaID, ejercicioID int32, in RutinaEjercicioUpdateInput) (RutinaEjercicio, error) {
	if rutinaID == 0 || ejercicioID == 0 {
		return RutinaEjercicio{}, fmt.Errorf("%w: rutina_id y ejercicio_id son obligatorios", ErrInvalid)
	}
	if err := validRutinaEjercicioNumeros(in.Orden, in.Series, in.Repeticiones, in.DescansoSegundos); err != nil {
		return RutinaEjercicio{}, err
	}

	row, err := s.store.UpdateRutinaEjercicio(ctx, db.UpdateRutinaEjercicioParams{
		RutinaID:         rutinaID,
		EjercicioID:      ejercicioID,
		Orden:            in.Orden,
		Series:           toNullInt32(in.Series),
		Repeticiones:     toNullInt32(in.Repeticiones),
		DescansoSegundos: toNullInt32(in.DescansoSegundos),
	})
	if err != nil {
		return RutinaEjercicio{}, mapNotFound(err)
	}
	return fromDBRutinaEjercicio(row), nil
}

func (s *RutinaEjercicioService) Delete(ctx context.Context, rutinaID, ejercicioID int32) error {
	if rutinaID == 0 || ejercicioID == 0 {
		return fmt.Errorf("%w: rutina_id y ejercicio_id son obligatorios", ErrInvalid)
	}
	return s.store.DeleteEjercicioDeRutina(ctx, db.DeleteEjercicioDeRutinaParams{
		RutinaID:    rutinaID,
		EjercicioID: ejercicioID,
	})
}

func validRutinaEjercicioNumeros(orden int32, series, repeticiones, descanso *int32) error {
	if orden <= 0 {
		return fmt.Errorf("%w: orden es obligatorio", ErrInvalid)
	}
	if !optionalPositive(series) || !optionalPositive(repeticiones) || !optionalPositive(descanso) {
		return fmt.Errorf("%w: series, repeticiones y descanso_segundos deben ser mayores a 0", ErrInvalid)
	}
	return nil
}

func optionalPositive(v *int32) bool {
	return v == nil || *v > 0
}
