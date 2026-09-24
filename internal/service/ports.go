package service

import (
	"context"

	db "PRACTICO_DOS/db/sqlc"
)

// EjercicioStore es el puerto de persistencia de ejercicios.
// *db.Queries lo satisface sin adaptador.
type EjercicioStore interface {
	CreateEjercicio(ctx context.Context, arg db.CreateEjercicioParams) (db.Ejercicio, error)
	// OPERACIONES CRUD
	GetEjercicio(ctx context.Context, id int32) (db.Ejercicio, error)
	ListEjercicios(ctx context.Context) ([]db.Ejercicio, error)
	UpdateEjercicio(ctx context.Context, arg db.UpdateEjercicioParams) error
	DeleteEjercicio(ctx context.Context, id int32) error
}

// AlumnoStore es el puerto de persistencia de alumnos.
// El alta usa CreateAlumno (usuario + alumno), no InsertAlumno suelto.
type AlumnoStore interface {
	CreateAlumno(ctx context.Context, arg db.CreateAlumnoParams) (db.GetAlumnoCompletoRow, error)
	GetAlumnoCompleto(ctx context.Context, id int32) (db.GetAlumnoCompletoRow, error)
	ListAlumnosCompletos(ctx context.Context) ([]db.ListAlumnosCompletosRow, error)
	UpdateAlumnoDetalles(ctx context.Context, arg db.UpdateAlumnoDetallesParams) error
	DeleteUsuario(ctx context.Context, id int32) error
}

// RutinaStore es el puerto de persistencia de rutinas.
// *db.Queries lo satisface sin adaptador.
type RutinaStore interface {
	CreateRutina(ctx context.Context, arg db.CreateRutinaParams) (db.Rutina, error)
	GetRutina(ctx context.Context, id int32) (db.Rutina, error)
	ListRutinas(ctx context.Context) ([]db.Rutina, error)
	UpdateRutina(ctx context.Context, arg db.UpdateRutinaParams) error
	DeleteRutina(ctx context.Context, id int32) error
}

// RutinaEjercicioStore es el puerto de la tabla intermedia rutina_ejercicio.
// El alta y la baja usan las queries que ya existían. El listado es por rutina.
type RutinaEjercicioStore interface {
	AddEjercicioARutina(ctx context.Context, arg db.AddEjercicioARutinaParams) error
	ListRutinaEjercicios(ctx context.Context, rutinaID int32) ([]db.RutinaEjercicio, error)
	UpdateRutinaEjercicio(ctx context.Context, arg db.UpdateRutinaEjercicioParams) (db.RutinaEjercicio, error)
	DeleteEjercicioDeRutina(ctx context.Context, arg db.DeleteEjercicioDeRutinaParams) error
}
