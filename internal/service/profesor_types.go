package service

import (
	"database/sql"

	db "PRACTICO_DOS/db/sqlc"
)

type Profesor struct {
	ID           int32   `json:"id"`
	Email        string  `json:"email"`
	Nombre       string  `json:"nombre"`
	Apellido     string  `json:"apellido"`
	Telefono     *string `json:"telefono"`
	Rol          string  `json:"rol"`
	CreatedAt    *string `json:"created_at"`
	Especialidad *string `json:"especialidad"`
}

// DEFINIMOS LOS DTOS
// PARA LAS OPERACIONES DE CREACION Y ACTUALIZACION DE PROFESORES
type ProfesorCreateInput struct {
	Email        string  `json:"email"`
	Contrasena   string  `json:"contrasena"`
	Nombre       string  `json:"nombre"`
	Apellido     string  `json:"apellido"`
	Telefono     *string `json:"telefono"`
	Especialidad *string `json:"especialidad"`
}

type ProfesorUpdateInput struct {
	Especialidad *string `json:"especialidad"`
}

func fromProfesorRow(id int32, email, nombre, apellido, rol string, telefono sql.NullString, createdAt sql.NullTime, especialidad sql.NullString) Profesor {
	return Profesor{
		ID:           id,
		Email:        email,
		Nombre:       nombre,
		Apellido:     apellido,
		Telefono:     fromNullString(telefono),
		Rol:          rol,
		CreatedAt:    fromNullTime(createdAt),
		Especialidad: fromNullString(especialidad),
	}
}

func fromDBProfesor(row db.GetProfesorCompletoRow) Profesor {
	return fromProfesorRow(row.ID, row.Email, row.Nombre, row.Apellido, row.Rol, row.Telefono, row.CreatedAt, row.Especialidad)
}

func fromDBProfesorList(row db.ListProfesoresCompletosRow) Profesor {
	return fromProfesorRow(row.ID, row.Email, row.Nombre, row.Apellido, row.Rol, row.Telefono, row.CreatedAt, row.Especialidad)
}
