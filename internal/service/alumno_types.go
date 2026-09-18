package service

import (
	"database/sql"
	"time"

	db "PRACTICO_DOS/db/sqlc"
)

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
