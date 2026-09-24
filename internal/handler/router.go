package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
)

func NewRouter(ejercicios *service.EjercicioService, alumnos *service.AlumnoService, rutinas *service.RutinaService, rutinas_ejercicio *service.RutinaEjercicioService, clases *service.ClaseService, clases_horario *service.ClaseHorarioService, staticDir string) http.Handler {
	ej := NewEjercicioHandler(ejercicios)
	al := NewAlumnoHandler(alumnos)
	ru := NewRutinaHandler(rutinas)
	re := NewRutinaEjercicioHandler(rutinas_ejercicio)
	cl := NewClaseHandler(clases)
	ch := NewClaseHorarioHandler(clases_horario)

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

	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux
}
