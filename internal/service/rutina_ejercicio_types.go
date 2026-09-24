package service

import (
	db "PRACTICO_DOS/db/sqlc"
)

type RutinaEjercicio struct {
	RutinaID         int32  `json:"rutina_id"`
	EjercicioID      int32  `json:"ejercicio_id"`
	Orden            int32  `json:"orden"`
	Series           *int32 `json:"series"`
	Repeticiones     *int32 `json:"repeticiones"`
	DescansoSegundos *int32 `json:"descanso_segundos"`
}

type RutinaEjercicioCreateInput struct {
	EjercicioID      int32  `json:"ejercicio_id"`
	Orden            int32  `json:"orden"`
	Series           *int32 `json:"series"`
	Repeticiones     *int32 `json:"repeticiones"`
	DescansoSegundos *int32 `json:"descanso_segundos"`
}

// No incluye ejercicio_id: cambiar el ejercicio es borrar la fila y crear otra.
type RutinaEjercicioUpdateInput struct {
	Orden            int32  `json:"orden"`
	Series           *int32 `json:"series"`
	Repeticiones     *int32 `json:"repeticiones"`
	DescansoSegundos *int32 `json:"descanso_segundos"`
}

func fromDBRutinaEjercicio(row db.RutinaEjercicio) RutinaEjercicio {
	return RutinaEjercicio{
		RutinaID:         row.RutinaID,
		EjercicioID:      row.EjercicioID,
		Orden:            row.Orden,
		Series:           fromNullInt32(row.Series),
		Repeticiones:     fromNullInt32(row.Repeticiones),
		DescansoSegundos: fromNullInt32(row.DescansoSegundos),
	}
}
