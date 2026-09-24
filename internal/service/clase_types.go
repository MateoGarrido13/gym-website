package service

import (
	db "PRACTICO_DOS/db/sqlc"
)

type Clase struct {
	ID         int32  `json:"id"`
	Nombre     string `json:"nombre"`
	MaxAlumnos int32  `json:"max_alumnos"`
	ProfesorID int32  `json:"profesor_id"`
}

type ClaseCreateInput struct {
	Nombre     string `json:"nombre"`
	MaxAlumnos int32  `json:"max_alumnos"`
	ProfesorID int32  `json:"profesor_id"`
}

// No incluye profesor_id: la query UpdateClase no lo cambia.
type ClaseUpdateInput struct {
	Nombre     string `json:"nombre"`
	MaxAlumnos int32  `json:"max_alumnos"`
}

func fromDBClase(row db.Clase) Clase {
	return Clase{
		ID:         row.ID,
		Nombre:     row.Nombre,
		MaxAlumnos: row.MaxAlumnos,
		ProfesorID: row.ProfesorID,
	}
}
