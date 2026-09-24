package service

import (
	"context"
	"fmt"

	db "PRACTICO_DOS/db/sqlc"
)

type ClaseHorarioService struct {
	store ClaseHorarioStore
}

func NewClaseHorarioService(store ClaseHorarioStore) *ClaseHorarioService {
	return &ClaseHorarioService{store: store}
}

func (s *ClaseHorarioService) Create(ctx context.Context, claseID int32, in ClaseHorarioCreateInput) (ClaseHorario, error) {
	if claseID == 0 {
		return ClaseHorario{}, fmt.Errorf("%w: clase_id es obligatorio", ErrInvalid)
	}

	dia, err := normalizeDiaSemana(in.DiaSemana)
	if err != nil {
		return ClaseHorario{}, err
	}
	hora, err := parseHoraInicio(in.HoraInicio)
	if err != nil {
		return ClaseHorario{}, err
	}

	row, err := s.store.CreateClaseHorario(ctx, db.CreateClaseHorarioParams{
		ClaseID:    claseID,
		DiaSemana:  dia,
		HoraInicio: hora,
	})
	if err != nil {
		return ClaseHorario{}, err
	}
	return fromDBClaseHorario(row), nil
}

func (s *ClaseHorarioService) List(ctx context.Context, claseID int32) ([]ClaseHorario, error) {
	if claseID == 0 {
		return nil, fmt.Errorf("%w: clase_id es obligatorio", ErrInvalid)
	}

	rows, err := s.store.ListHorariosDeClase(ctx, claseID)
	if err != nil {
		return nil, err
	}
	out := make([]ClaseHorario, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBClaseHorario(row))
	}
	return out, nil
}

func (s *ClaseHorarioService) Get(ctx context.Context, id int32) (ClaseHorario, error) {
	row, err := s.store.GetClaseHorario(ctx, id)
	if err != nil {
		return ClaseHorario{}, mapNotFound(err)
	}
	return fromDBClaseHorario(row), nil
}

func (s *ClaseHorarioService) Delete(ctx context.Context, id int32) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteClaseHorario(ctx, id)
}
