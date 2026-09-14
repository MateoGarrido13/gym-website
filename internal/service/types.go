package service

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	db "PRACTICO_DOS/db/sqlc"
)

var (
	ErrInvalid  = errors.New("datos inválidos")
	ErrNotFound = errors.New("no encontrado")
)

type Ejercicio struct {
	ID            int32   `json:"id"`
	Nombre        string  `json:"nombre"`
	Descripcion   *string `json:"descripcion"`
	GrupoMuscular *string `json:"grupo_muscular"`
}

type EjercicioInput struct {
	Nombre        string  `json:"nombre"`
	Descripcion   *string `json:"descripcion"`
	GrupoMuscular *string `json:"grupo_muscular"`
}

type Alumno struct {
	ID               int32   `json:"id"`
	Email            string  `json:"email"`
	Nombre           string  `json:"nombre"`
	Apellido         string  `json:"apellido"`
	Telefono         *string `json:"telefono"`
	Rol              string  `json:"rol"`
	CreatedAt        *string `json:"created_at"`
	FechaInscripcion string  `json:"fecha_inscripcion"`
	FechaVto         *string `json:"fecha_vto"`
	TipoPlan         *string `json:"tipo_plan"`
	RutinaID         *int32  `json:"rutina_id"`
}

type AlumnoCreateInput struct {
	Email            string  `json:"email"`
	Contrasena       string  `json:"contrasena"`
	Nombre           string  `json:"nombre"`
	Apellido         string  `json:"apellido"`
	Telefono         *string `json:"telefono"`
	FechaInscripcion string  `json:"fecha_inscripcion"`
	FechaVto         *string `json:"fecha_vto"`
	TipoPlan         *string `json:"tipo_plan"`
	RutinaID         *int32  `json:"rutina_id"`
}

// DEFINIMOS LOS DTOS 
// PARA LAS OPERACIONES DE ACTUALIZACION Y CREACION DE ALUMNOS
type AlumnoUpdateInput struct {
	FechaVto *string `json:"fecha_vto"`
	TipoPlan *string `json:"tipo_plan"`
	RutinaID *int32  `json:"rutina_id"`
}

// FUNCIONES DE CONVERSION DE LOS DATOS DE LA BASE DE DATOS A LOS DTOS
func fromDBEjercicio(e db.Ejercicio) Ejercicio {
	return Ejercicio{
		ID:            e.ID,
		Nombre:        e.Nombre,
		Descripcion:   fromNullString(e.Descripcion),
		GrupoMuscular: fromNullString(e.GrupoMuscular),
	}
}

// FUNCIONES DE CONVERSION DE LOS DATOS DE LA BASE DE DATOS A LOS DTOS
func fromAlumnoRow(id int32, email, nombre, apellido, rol string, telefono sql.NullString, createdAt sql.NullTime, fechaInscripcion time.Time, fechaVto sql.NullTime, tipoPlan sql.NullString, rutinaID sql.NullInt32) Alumno {
	return Alumno{
		ID:               id,
		Email:            email,
		Nombre:           nombre,
		Apellido:         apellido,
		Telefono:         fromNullString(telefono),
		Rol:              rol,
		CreatedAt:        fromNullTime(createdAt),
		FechaInscripcion: fechaInscripcion.UTC().Format(time.DateOnly),
		FechaVto:         fromNullTimeDate(fechaVto),
		TipoPlan:         fromNullString(tipoPlan),
		RutinaID:         fromNullInt32(rutinaID),
	}
}

func fromDBAlumno(row db.GetAlumnoCompletoRow) Alumno {
	return fromAlumnoRow(row.ID, row.Email, row.Nombre, row.Apellido, row.Rol, row.Telefono, row.CreatedAt, row.FechaInscripcion, row.FechaVto, row.TipoPlan, row.RutinaID)
}

func fromDBAlumnoList(row db.ListAlumnosCompletosRow) Alumno {
	return fromAlumnoRow(row.ID, row.Email, row.Nombre, row.Apellido, row.Rol, row.Telefono, row.CreatedAt, row.FechaInscripcion, row.FechaVto, row.TipoPlan, row.RutinaID)
}


func fromNullString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func fromNullInt32(v sql.NullInt32) *int32 {
	if !v.Valid {
		return nil
	}
	n := v.Int32
	return &n
}

func fromNullTime(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	s := v.Time.UTC().Format(time.RFC3339)
	return &s
}

func fromNullTimeDate(v sql.NullTime) *string {
	if !v.Valid {
		return nil
	}
	s := v.Time.UTC().Format(time.DateOnly)
	return &s
}

func toNullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func toNullInt32(v *int32) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *v, Valid: true}
}

func parseDate(s string) (time.Time, error) {
	return time.Parse(time.DateOnly, strings.TrimSpace(s))
}

func toNullDate(v *string) (sql.NullTime, error) {
	if v == nil || strings.TrimSpace(*v) == "" {
		return sql.NullTime{}, nil
	}
	t, err := parseDate(*v)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}
