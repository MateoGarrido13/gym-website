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
	ProfesorID      int32  `json:"profesor_id"`
}

func fromDBRutina(row db.Rutina) Rutina {
	return Rutina{
		ID:              row.ID,
		Nombre:          row.Nombre,
		DuracionSemanas: row.DuracionSemanas,
		ProfesorID:      row.ProfesorID,
	}
}

// RutinaAlumno es la rutina asignada a un alumno, lista para mostrarla: con el
// nombre del profesor y los ejercicios en orden.
type RutinaAlumno struct {
	ID              int32               `json:"id"`
	Nombre          string              `json:"nombre"`
	DuracionSemanas int32               `json:"duracion_semanas"`
	Profesor        string              `json:"profesor"`
	Ejercicios      []EjercicioDeRutina `json:"ejercicios"`
}

type EjercicioDeRutina struct {
	Orden            int32   `json:"orden"`
	Nombre           string  `json:"nombre"`
	Descripcion      *string `json:"descripcion"`
	GrupoMuscular    *string `json:"grupo_muscular"`
	Series           *int32  `json:"series"`
	Repeticiones     *int32  `json:"repeticiones"`
	DescansoSegundos *int32  `json:"descanso_segundos"`
}

func fromDBEjercicioDeRutina(row db.ListEjerciciosDetalleDeRutinaRow) EjercicioDeRutina {
	return EjercicioDeRutina{
		Orden:            row.Orden,
		Nombre:           row.Nombre,
		Descripcion:      fromNullString(row.Descripcion),
		GrupoMuscular:    fromNullString(row.GrupoMuscular),
		Series:           fromNullInt32(row.Series),
		Repeticiones:     fromNullInt32(row.Repeticiones),
		DescansoSegundos: fromNullInt32(row.DescansoSegundos),
	}
}
