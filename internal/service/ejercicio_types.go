package service

import (
	db "PRACTICO_DOS/db/sqlc"
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

func fromDBEjercicio(e db.Ejercicio) Ejercicio {
	return Ejercicio{
		ID:            e.ID,
		Nombre:        e.Nombre,
		Descripcion:   fromNullString(e.Descripcion),
		GrupoMuscular: fromNullString(e.GrupoMuscular),
	}
}
