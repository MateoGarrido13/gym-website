package service

import (
	"context"
	"fmt"
	"strings"

	db "PRACTICO_DOS/db/sqlc"
)

type AlumnoService struct {
	store AlumnoStore
}

func NewAlumnoService(store AlumnoStore) *AlumnoService {
	return &AlumnoService{store: store}
}

func (s *AlumnoService) Create(ctx context.Context, in AlumnoCreateInput) (Alumno, error) {
	email := strings.TrimSpace(in.Email)
	nombre := strings.TrimSpace(in.Nombre)
	apellido := strings.TrimSpace(in.Apellido)
	if email == "" || in.Contrasena == "" || nombre == "" || apellido == "" {
		return Alumno{}, fmt.Errorf("%w: email, contraseña, nombre y apellido son obligatorios", ErrInvalid)
	}

	fechaInscripcion, err := parseDate(in.FechaInscripcion)
	if err != nil {
		return Alumno{}, fmt.Errorf("%w: fecha_inscripcion debe ser YYYY-MM-DD", ErrInvalid)
	}
	fechaVto, err := toNullDate(in.FechaVto)
	if err != nil {
		return Alumno{}, fmt.Errorf("%w: fecha_vto debe ser YYYY-MM-DD", ErrInvalid)
	}

	row, err := s.store.CreateAlumno(ctx, db.CreateAlumnoParams{
		Email:            email,
		Contrasena:       in.Contrasena,
		Nombre:           nombre,
		Apellido:         apellido,
		Telefono:         toNullString(in.Telefono),
		FechaInscripcion: fechaInscripcion,
		FechaVto:         fechaVto,
		TipoPlan:         toNullString(in.TipoPlan),
		RutinaID:         toNullInt32(in.RutinaID),
	})
	if err != nil {
		return Alumno{}, err
	}
	return fromDBAlumno(row), nil
}

func (s *AlumnoService) Get(ctx context.Context, id int32) (Alumno, error) {
	row, err := s.store.GetAlumnoCompleto(ctx, id)
	if err != nil {
		return Alumno{}, mapNotFound(err)
	}
	return fromDBAlumno(row), nil
}

func (s *AlumnoService) List(ctx context.Context) ([]Alumno, error) {
	rows, err := s.store.ListAlumnosCompletos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]Alumno, 0, len(rows))
	for _, row := range rows {
		out = append(out, fromDBAlumnoList(row))
	}
	return out, nil
}

func (s *AlumnoService) Update(ctx context.Context, id int32, in AlumnoUpdateInput) (Alumno, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return Alumno{}, err
	}

	fechaVto, err := toNullDate(in.FechaVto)
	if err != nil {
		return Alumno{}, fmt.Errorf("%w: fecha_vto debe ser YYYY-MM-DD", ErrInvalid)
	}

	err = s.store.UpdateAlumnoDetalles(ctx, db.UpdateAlumnoDetallesParams{
		UsuarioID: id,
		FechaVto:  fechaVto,
		TipoPlan:  toNullString(in.TipoPlan),
		RutinaID:  toNullInt32(in.RutinaID),
	})
	if err != nil {
		return Alumno{}, err
	}
	return s.Get(ctx, id)
}

func (s *AlumnoService) Delete(ctx context.Context, id int32) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteUsuario(ctx, id)
}
