package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
)

func NewRouter(ejercicios *service.EjercicioService, alumnos *service.AlumnoService, rutinas *service.RutinaService, rutinas_ejercicio *service.RutinaEjercicioService, staticDir string) http.Handler {
	ej := NewEjercicioHandler(ejercicios)
	al := NewAlumnoHandler(alumnos)
	ru := NewRutinaHandler(rutinas)
	re := NewRutinaEjercicioHandler(rutinas_ejercicio)

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

	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux
}
