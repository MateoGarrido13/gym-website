package service

import (
	"time"

	db "PRACTICO_DOS/db/sqlc"
)

// SesionDuracion es el tiempo de vida fijo de una sesión (sin renovación).
const SesionDuracion = 24 * time.Hour

// UsuarioSesion son los datos del usuario autenticado. No incluye la contraseña.
type UsuarioSesion struct {
	ID       int32  `json:"id"`
	Email    string `json:"email"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Rol      string `json:"rol"`
}

// SesionIniciada es lo que devuelve Login: el token va en la cookie del
// navegador; Expira sirve para fijar el vencimiento de la cookie.
type SesionIniciada struct {
	Token   string
	Expira  time.Time
	Usuario UsuarioSesion
}

func fromDBUsuarioSesion(row db.GetSesionUsuarioRow) UsuarioSesion {
	return UsuarioSesion{ID: row.ID, Email: row.Email, Nombre: row.Nombre, Apellido: row.Apellido, Rol: row.Rol}
}
