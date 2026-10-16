package service

import (
	"context"
	"fmt"
	"strings"

	db "PRACTICO_DOS/db/sqlc"
)

type ProfesorService struct {
	store ProfesorStore
}

func NewProfesorService(store ProfesorStore) *ProfesorService {
	return &ProfesorService{store: store}
}

func (s *ProfesorService) Create(ctx context.Context, in ProfesorCreateInput) (Profesor, error) {
	email := strings.TrimSpace(in.Email)
	nombre := strings.TrimSpace(in.Nombre)
	apellido := strings.TrimSpace(in.Apellido)
	if email == "" || in.Contrasena == "" || nombre == "" || apellido == "" {
		return Profesor{}, fmt.Errorf("%w: email, contraseña, nombre y apellido son obligatorios", ErrInvalid)
	}

	row, err := s.store.CreateProfesor(ctx, db.CreateProfesorParams{
		Email:        email,
		Contrasena:   in.Contrasena,
		Nombre:       nombre,
		Apellido:     apellido,
		Telefono:     toNullString(in.Telefono),
		Especialidad: toNullString(in.Especialidad),
	})
	if err != nil {
		return Profesor{}, err
	}
	return fromDBProfesor(row), nil
}

func (s *ProfesorService) Get(ctx context.Context, id int32) (Profesor, error) {
	row, err := s.store.GetProfesorCompleto(ctx, id)
	if err != nil {
		return Profesor{}, mapNotFound(err)
	}
	return fromDBProfesor(row), nil
}

func (s *ProfesorService) List(ctx context.Context) ([]Profesor, error) {
	rows, err := s.store.ListProfesoresCompletos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Profesor, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBProfesorList(row))
	}
	return out, nil
}

func (s *ProfesorService) Update(ctx context.Context, id int32, in ProfesorUpdateInput) (Profesor, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return Profesor{}, err
	}

	err := s.store.UpdateProfesor(ctx, db.UpdateProfesorParams{
		UsuarioID:    id,
		Especialidad: toNullString(in.Especialidad),
	})
	if err != nil {
		return Profesor{}, err
	}
	return s.Get(ctx, id)
}

func (s *ProfesorService) Delete(ctx context.Context, id int32) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteUsuario(ctx, id)
}