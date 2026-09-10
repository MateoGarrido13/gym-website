package db

import (
	"context"
	"database/sql"
	"time"
)

// CreateAlumnoParams agrupa los datos del usuario padre y del alumno.
type CreateAlumnoParams struct {
	Email            string         `json:"email"`
	Contrasena       string         `json:"contrasena"`
	Nombre           string         `json:"nombre"`
	Apellido         string         `json:"apellido"`
	Telefono         sql.NullString `json:"telefono"`
	FechaInscripcion time.Time      `json:"fecha_inscripcion"`
	FechaVto         sql.NullTime   `json:"fecha_vto"`
	TipoPlan         sql.NullString `json:"tipo_plan"`
	RutinaID         sql.NullInt32  `json:"rutina_id"`
}

// CreateAlumno inserta primero el usuario padre y después la fila de alumno.
func (q *Queries) CreateAlumno(ctx context.Context, arg CreateAlumnoParams) (GetAlumnoCompletoRow, error) {
	usuario, err := q.CreateUsuario(ctx, CreateUsuarioParams{ //
		Email:      arg.Email,
		Contrasena: arg.Contrasena,
		Nombre:     arg.Nombre,
		Apellido:   arg.Apellido,
		Telefono:   arg.Telefono,
		Rol:        "alumno",
	})
	if err != nil {
		return GetAlumnoCompletoRow{}, err
	}
	
	err = q.InsertAlumno(ctx, InsertAlumnoParams{
		UsuarioID:        usuario.ID,
		FechaInscripcion: arg.FechaInscripcion,
		FechaVto:         arg.FechaVto,
		TipoPlan:         arg.TipoPlan,
		RutinaID:         arg.RutinaID,
	})
	if err != nil {
		_ = q.DeleteUsuario(ctx, usuario.ID)
		return GetAlumnoCompletoRow{}, err
	}

	return q.GetAlumnoCompleto(ctx, usuario.ID)
}
