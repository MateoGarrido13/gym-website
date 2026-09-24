package service

import (
	"context"
	"fmt"
	"strings"

	db "PRACTICO_DOS/db/sqlc"
)

type ClaseService struct {
	store ClaseStore
}

func NewClaseService(store ClaseStore) *ClaseService {
	return &ClaseService{store: store}
}

func (s *ClaseService) Create(ctx context.Context, in ClaseCreateInput) (Clase, error) {
	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" || in.MaxAlumnos <= 0 || in.ProfesorID == 0 {
		return Clase{}, fmt.Errorf("%w: nombre, max_alumnos y profesor_id son obligatorios", ErrInvalid)
	}

	row, err := s.store.CreateClase(ctx, db.CreateClaseParams{
		Nombre:     nombre,
		MaxAlumnos: in.MaxAlumnos,
		ProfesorID: in.ProfesorID,
	})
	if err != nil {
		return Clase{}, err
	}
	return fromDBClase(row), nil
}

func (s *ClaseService) Get(ctx context.Context, id int32) (Clase, error) {
	row, err := s.store.GetClase(ctx, id)
	if err != nil {
		return Clase{}, mapNotFound(err)
	}
	return fromDBClase(row), nil
}

func (s *ClaseService) List(ctx context.Context) ([]Clase, error) {
	rows, err := s.store.ListClases(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Clase, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBClase(row))
	}
	return out, nil
}

func (s *ClaseService) Update(ctx context.Context, id int32, in ClaseUpdateInput) (Clase, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return Clase{}, err
	}

	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" || in.MaxAlumnos <= 0 {
		return Clase{}, fmt.Errorf("%w: nombre y max_alumnos son obligatorios", ErrInvalid)
	}

	err := s.store.UpdateClase(ctx, db.UpdateClaseParams{
		ID:         id,
		Nombre:     nombre,
		MaxAlumnos: in.MaxAlumnos,
	})
	if err != nil {
		return Clase{}, err
	}
	return s.Get(ctx, id)
}

func (s *ClaseService) Delete(ctx context.Context, id int32) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteClase(ctx, id)
}
