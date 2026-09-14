package db

import (
	"context"
	"database/sql"
)

// CreateProfesorParams agrupa los datos del usuario padre y del profesor.
type CreateProfesorParams struct {
	Email        string         `json:"email"`
	Contrasena   string         `json:"contrasena"`
	Nombre       string         `json:"nombre"`
	Apellido     string         `json:"apellido"`
	Telefono     sql.NullString `json:"telefono"`
	Especialidad sql.NullString `json:"especialidad"`
}

// CreateProfesor inserta primero el usuario padre y después la fila de profesor.
func (q *Queries) CreateProfesor(ctx context.Context, arg CreateProfesorParams) (GetProfesorCompletoRow, error) {
	usuario, err := q.CreateUsuario(ctx, CreateUsuarioParams{
		Email:      arg.Email,
		Contrasena: arg.Contrasena,
		Nombre:     arg.Nombre,
		Apellido:   arg.Apellido,
		Telefono:   arg.Telefono,
		Rol:        "profesor",
	})
	if err != nil {
		return GetProfesorCompletoRow{}, err
	}

	err = q.InsertProfesor(ctx, InsertProfesorParams{
		UsuarioID:    usuario.ID,
		Especialidad: arg.Especialidad,
	})
	if err != nil {
		_ = q.DeleteUsuario(ctx, usuario.ID)
		return GetProfesorCompletoRow{}, err
	}

	return q.GetProfesorCompleto(ctx, usuario.ID)
}
