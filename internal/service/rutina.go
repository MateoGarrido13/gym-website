package service

import (
	"context"
	"fmt"
	"strings"

	db "PRACTICO_DOS/db/sqlc"
)

type RutinaService struct {
	store RutinaStore
}

func NewRutinaService(store RutinaStore) *RutinaService {
	return &RutinaService{store: store}
}

func (s *RutinaService) Create(ctx context.Context, in RutinaCreateInput) (Rutina, error) {
	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" || in.DuracionSemanas == 0 || in.ProfesorID == 0 {
		return Rutina{}, fmt.Errorf("%w: nombre, duracion_semanas y profesor_id son obligatorios", ErrInvalid)
	}

	row, err := s.store.CreateRutina(ctx, db.CreateRutinaParams{
		Nombre:          nombre,
		DuracionSemanas: in.DuracionSemanas,
		ProfesorID:      in.ProfesorID,
	})
	if err != nil {
		return Rutina{}, err
	}
	return fromDBRutina(row), nil
}

func (s *RutinaService) Get(ctx context.Context, id int32) (Rutina, error) {
	row, err := s.store.GetRutina(ctx, id)
	if err != nil {
		return Rutina{}, mapNotFound(err)
	}
	return fromDBRutina(row), nil
}

func (s *RutinaService) List(ctx context.Context) ([]Rutina, error) {
	rows, err := s.store.ListRutinas(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Rutina, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBRutina(row))
	}
	return out, nil
}

func (s *RutinaService) Update(ctx context.Context, id int32, in RutinaUpdateInput) (Rutina, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return Rutina{}, err
	}

	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" || in.DuracionSemanas == 0 {
		return Rutina{}, fmt.Errorf("%w: nombre y duracion_semanas son obligatorios", ErrInvalid)
	}

	err := s.store.UpdateRutina(ctx, db.UpdateRutinaParams{
		ID:              id,
		Nombre:          nombre,
		DuracionSemanas: in.DuracionSemanas,
	})
	if err != nil {
		return Rutina{}, err
	}
	return s.Get(ctx, id)
}

func (s *RutinaService) Delete(ctx context.Context, id int32) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteRutina(ctx, id)
}
