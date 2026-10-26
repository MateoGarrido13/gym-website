package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/internal/web"
)

func NewRouter(ejercicios *service.EjercicioService, alumnos *service.AlumnoService, profesores *service.ProfesorService, rutinas *service.RutinaService, rutinas_ejercicio *service.RutinaEjercicioService, clases *service.ClaseService, clases_horario *service.ClaseHorarioService, inscripciones *service.InscripcionService, staticDir string) http.Handler {
	ej := NewEjercicioHandler(ejercicios)
	al := NewAlumnoHandler(alumnos)
	pr := NewProfesorHandler(profesores)
	ru := NewRutinaHandler(rutinas)
	re := NewRutinaEjercicioHandler(rutinas_ejercicio)
	cl := NewClaseHandler(clases)
	ch := NewClaseHorarioHandler(clases_horario)
	ins := NewInscripcionHandler(inscripciones)
	pages := web.NewHandler(ejercicios, alumnos, profesores, rutinas, rutinas_ejercicio, clases)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/ejercicios", ej.Create)
	mux.HandleFunc("GET /api/ejercicios", ej.List)
	mux.HandleFunc("GET /api/ejercicios/{id}", ej.Get)
	mux.HandleFunc("PUT /api/ejercicios/{id}", ej.Update)
	mux.HandleFunc("DELETE /api/ejercicios/{id}", ej.Delete)

	mux.HandleFunc("POST /api/alumnos", al.Create)
	mux.HandleFunc("GET /api/alumnos", al.List)
	mux.HandleFunc("GET /api/alumnos/{id}", al.Get)
	mux.HandleFunc("PUT /api/alumnos/{id}", al.Update)
	mux.HandleFunc("DELETE /api/alumnos/{id}", al.Delete)

	mux.HandleFunc("POST /api/profesores", pr.Create)
	mux.HandleFunc("GET /api/profesores", pr.List)
	mux.HandleFunc("GET /api/profesores/{id}", pr.Get)
	mux.HandleFunc("PUT /api/profesores/{id}", pr.Update)
	mux.HandleFunc("DELETE /api/profesores/{id}", pr.Delete)

	mux.HandleFunc("POST /api/rutinas", ru.Create)
	mux.HandleFunc("GET /api/rutinas", ru.List)
	mux.HandleFunc("GET /api/rutinas/{id}", ru.Get)
	mux.HandleFunc("PUT /api/rutinas/{id}", ru.Update)
	mux.HandleFunc("DELETE /api/rutinas/{id}", ru.Delete)

	mux.HandleFunc("POST /api/rutinas/{id}/ejercicios", re.Create)
	mux.HandleFunc("GET /api/rutinas/{id}/ejercicios", re.List)
	mux.HandleFunc("PUT /api/rutinas/{id}/ejercicios/{ejercicio_id}", re.Update)
	mux.HandleFunc("DELETE /api/rutinas/{id}/ejercicios/{ejercicio_id}", re.Delete)

	mux.HandleFunc("POST /api/clases", cl.Create)
	mux.HandleFunc("GET /api/clases", cl.List)
	mux.HandleFunc("GET /api/clases/{id}", cl.Get)
	mux.HandleFunc("PUT /api/clases/{id}", cl.Update)
	mux.HandleFunc("DELETE /api/clases/{id}", cl.Delete)

	mux.HandleFunc("POST /api/clases/{id}/horarios", ch.Create)
	mux.HandleFunc("GET /api/clases/{id}/horarios", ch.List)
	mux.HandleFunc("GET /api/horarios/{id}", ch.Get)
	mux.HandleFunc("DELETE /api/horarios/{id}", ch.Delete)

	mux.HandleFunc("POST /api/horarios/{id}/inscripciones", ins.Create)
	mux.HandleFunc("GET /api/horarios/{id}/inscripciones", ins.List)
	mux.HandleFunc("DELETE /api/horarios/{id}/inscripciones/{alumno_id}", ins.Delete)

	mux.HandleFunc("GET /", pages.Alumnos)
	mux.HandleFunc("GET /profesores", pages.Profesores)
	mux.HandleFunc("POST /alumnos", pages.CreateAlumno)
	mux.HandleFunc("POST /profesores", pages.CreateProfesor)
	mux.HandleFunc("GET /alumnos/{id}/editar", pages.EditAlumno)
	mux.HandleFunc("POST /alumnos/{id}", pages.UpdateAlumno)
	mux.HandleFunc("POST /alumnos/{id}/eliminar", pages.DeleteAlumno)
	mux.HandleFunc("GET /profesores/{id}/editar", pages.EditProfesor)
	mux.HandleFunc("POST /profesores/{id}", pages.UpdateProfesor)
	mux.HandleFunc("POST /profesores/{id}/eliminar", pages.DeleteProfesor)
	mux.HandleFunc("GET /ejercicios", pages.Ejercicios)
	mux.HandleFunc("POST /ejercicios", pages.CreateEjercicio)
	mux.HandleFunc("GET /ejercicios/{id}/editar", pages.EditEjercicio)
	mux.HandleFunc("POST /ejercicios/{id}", pages.UpdateEjercicio)
	mux.HandleFunc("POST /ejercicios/{id}/eliminar", pages.DeleteEjercicio)
	mux.HandleFunc("GET /rutinas", pages.Rutinas)
	mux.HandleFunc("POST /rutinas", pages.CreateRutina)
	mux.HandleFunc("GET /rutinas/{id}/editar", pages.EditRutina)
	mux.HandleFunc("POST /rutinas/{id}", pages.UpdateRutina)
	mux.HandleFunc("POST /rutinas/{id}/eliminar", pages.DeleteRutina)
	mux.HandleFunc("GET /clases", pages.Clases)
	mux.HandleFunc("POST /clases", pages.CreateClase)
	mux.HandleFunc("GET /clases/{id}/editar", pages.EditClase)
	mux.HandleFunc("POST /clases/{id}", pages.UpdateClase)
	mux.HandleFunc("POST /clases/{id}/eliminar", pages.DeleteClase)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	return mux
}
