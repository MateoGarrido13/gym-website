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

// ClaseStore es el puerto de persistencia de clases.
// *db.Queries lo satisface sin adaptador.
type ClaseStore interface {
	CreateClase(ctx context.Context, arg db.CreateClaseParams) (db.Clase, error)
	GetClase(ctx context.Context, id int32) (db.Clase, error)
	ListClases(ctx context.Context) ([]db.Clase, error)
	UpdateClase(ctx context.Context, arg db.UpdateClaseParams) error
	DeleteClase(ctx context.Context, id int32) error
}

// ClaseHorarioStore es el puerto de los horarios de una clase.
// El listado es por clase. Get y Delete usan el id del horario.
type ClaseHorarioStore interface {
	CreateClaseHorario(ctx context.Context, arg db.CreateClaseHorarioParams) (db.ClaseHorario, error)
	GetClaseHorario(ctx context.Context, id int32) (db.ClaseHorario, error)
	ListHorariosDeClase(ctx context.Context, claseID int32) ([]db.ClaseHorario, error)
	DeleteClaseHorario(ctx context.Context, id int32) error
}

// InscripcionStore es el puerto de las inscripciones a un horario.
// El cupo se lee con GetCupoDeHorario; la comparación la hace el servicio.
type InscripcionStore interface {
	CreateInscripcion(ctx context.Context, arg db.CreateInscripcionParams) (db.InscripcionClase, error)
	ListInscriptos(ctx context.Context, claseHorarioID int32) ([]db.ListInscriptosRow, error)
	DeleteInscripcion(ctx context.Context, arg db.DeleteInscripcionParams) error
	GetCupoDeHorario(ctx context.Context, id int32) (db.GetCupoDeHorarioRow, error)
}
