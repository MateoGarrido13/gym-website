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
	if in.GrupoMuscular == nil || strings.TrimSpace(*in.GrupoMuscular) == "" {
		return Ejercicio{}, fmt.Errorf("%w: grupo muscular es obligatorio", ErrInvalid)
	}
	if err := s.ensureUniqueNombreGrupo(ctx, nombre, strings.TrimSpace(*in.GrupoMuscular), 0); err != nil {
		return Ejercicio{}, err
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
// Si existe, devuelve el ejercicio local.
// Si no existe, lo crea usando los datos de wger y devuelve el ejercicio local.
func (s *EjercicioService) FindOrCreateByNombreGrupo(ctx context.Context, nombre, grupoMuscular string, descripcion *string) (Ejercicio, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return Ejercicio{}, fmt.Errorf("%w: nombre es obligatorio", ErrInvalid)
	}
	if grupoMuscular == "" {
		grupoMuscular = "Sin categoría"
	}

	// Buscar coincidencia existente (case-insensitive).
	existentes, err := s.List(ctx)
	if err != nil {
		return Ejercicio{}, err
	}
	for _, ex := range existentes {
		if strings.EqualFold(nombre, ex.Nombre) &&
			strings.EqualFold(grupoMuscular, stringValueOrEmpty(ex.GrupoMuscular)) {
			return ex, nil
		}
	}

	// Crear el ejercicio local para poder respetar la FK.
	return s.Create(ctx, EjercicioInput{
		Nombre:        nombre,
		Descripcion:   descripcion,
		GrupoMuscular: &grupoMuscular,
	})
}

func (s *EjercicioService) Update(ctx context.Context, id int32, in EjercicioInput) (Ejercicio, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return Ejercicio{}, err
	}

	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" {
		return Ejercicio{}, fmt.Errorf("%w: nombre es obligatorio", ErrInvalid)
	}
	if in.GrupoMuscular == nil || strings.TrimSpace(*in.GrupoMuscular) == "" {
		return Ejercicio{}, fmt.Errorf("%w: grupo muscular es obligatorio", ErrInvalid)
	}
	if err := s.ensureUniqueNombreGrupo(ctx, nombre, strings.TrimSpace(*in.GrupoMuscular), id); err != nil {
		return Ejercicio{}, err
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

func (s *EjercicioService) ensureUniqueNombreGrupo(ctx context.Context, nombre, grupo string, ignoreID int32) error {
	rows, err := s.List(ctx)
	if err != nil {
		return err
	}
	for _, row := range rows {
		if row.ID == ignoreID || row.GrupoMuscular == nil {
			continue
		}
		sameName := strings.EqualFold(nombre, strings.TrimSpace(row.Nombre))
		sameGroup := strings.EqualFold(grupo, strings.TrimSpace(*row.GrupoMuscular))
		if sameName && sameGroup {
			return fmt.Errorf("%w: ya existe un ejercicio con el mismo nombre y grupo muscular", ErrInvalid)
		}
	}
	return nil
}

func (s *EjercicioService) Delete(ctx context.Context, id int32) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteEjercicio(ctx, id)
}

func stringValueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
