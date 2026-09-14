package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	db "PRACTICO_DOS/db/sqlc"
)

type EjercicioService struct {
	store EjercicioStore
}

func NewEjercicioService(store EjercicioStore) *EjercicioService {
	return &EjercicioService{store: store}
}

func (s *EjercicioService) Create(ctx context.Context, in EjercicioInput) (Ejercicio, error) {
	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" {
		return Ejercicio{}, fmt.Errorf("%w: nombre es obligatorio", ErrInvalid)
	}

	row, err := s.store.CreateEjercicio(ctx, db.CreateEjercicioParams{
		Nombre:        nombre,
		Descripcion:   toNullString(in.Descripcion),
		GrupoMuscular: toNullString(in.GrupoMuscular),
	})
	if err != nil {
		return Ejercicio{}, err
	}
	return fromDBEjercicio(row), nil
}

func (s *EjercicioService) Get(ctx context.Context, id int32) (Ejercicio, error) {
	row, err := s.store.GetEjercicio(ctx, id)
	if err != nil {
		return Ejercicio{}, mapNotFound(err)
	}
	return fromDBEjercicio(row), nil
}

func (s *EjercicioService) List(ctx context.Context) ([]Ejercicio, error) {
	rows, err := s.store.ListEjercicios(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Ejercicio, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBEjercicio(row))
	}
	return out, nil
}

func (s *EjercicioService) Update(ctx context.Context, id int32, in EjercicioInput) (Ejercicio, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return Ejercicio{}, err
	}

	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" {
		return Ejercicio{}, fmt.Errorf("%w: nombre es obligatorio", ErrInvalid)
	}

	err := s.store.UpdateEjercicio(ctx, db.UpdateEjercicioParams{
		ID:            id,
		Nombre:        nombre,
		Descripcion:   toNullString(in.Descripcion),
		GrupoMuscular: toNullString(in.GrupoMuscular),
	})
	if err != nil {
		return Ejercicio{}, err
	}
	return s.Get(ctx, id)
}

func (s *EjercicioService) Delete(ctx context.Context, id int32) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteEjercicio(ctx, id)
}

func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
