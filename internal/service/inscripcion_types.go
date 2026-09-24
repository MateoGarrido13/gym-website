package service

import (
	db "PRACTICO_DOS/db/sqlc"
)

type Inscripcion struct {
	ClaseHorarioID   int32   `json:"clase_horario_id"`
	AlumnoID         int32   `json:"alumno_id"`
	FechaInscripcion *string `json:"fecha_inscripcion"`
}

type InscripcionCreateInput struct {
	AlumnoID int32 `json:"alumno_id"`
}

type Inscripto struct {
	ID       int32  `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

func fromDBInscripcion(row db.InscripcionClase) Inscripcion {
	return Inscripcion{
		ClaseHorarioID:   row.ClaseHorarioID,
		AlumnoID:         row.AlumnoID,
		FechaInscripcion: fromNullTime(row.FechaInscripcion),
	}
}

func fromDBInscripto(row db.ListInscriptosRow) Inscripto {
	return Inscripto{
		ID:       row.AlumnoID,
		Nombre:   row.Nombre,
		Apellido: row.Apellido,
	}
}
