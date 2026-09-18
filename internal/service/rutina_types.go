package service

import (
	db "PRACTICO_DOS/db/sqlc"
)

type Rutina struct {
	ID              int32  `json:"id"`
	Nombre          string `json:"nombre"`
	DuracionSemanas int32  `json:"duracion_semanas"`
	ProfesorID      int32  `json:"profesor_id"`
}

type RutinaCreateInput struct {
	Nombre          string `json:"nombre"`
	DuracionSemanas int32  `json:"duracion_semanas"`
	ProfesorID      int32  `json:"profesor_id"`
}

type RutinaUpdateInput struct {
	Nombre          string `json:"nombre"`
	DuracionSemanas int32  `json:"duracion_semanas"`
}

func fromDBRutina(row db.Rutina) Rutina {
	return Rutina{
		ID:              row.ID,
		Nombre:          row.Nombre,
		DuracionSemanas: row.DuracionSemanas,
		ProfesorID:      row.ProfesorID,
	}
}