package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/internal/web"
)
//hay endpoints que ya no se consultan por el frontend, pero seran útiles para demos con curl/ postman
func NewRouter(ejercicios *service.EjercicioService, alumnos *service.AlumnoService, profesores *service.ProfesorService, rutinas *service.RutinaService, rutinas_ejercicio *service.RutinaEjercicioService, clases *service.ClaseService, clases_horario *service.ClaseHorarioService, inscripciones *service.InscripcionService, auth *service.AuthService, staticDir string) http.Handler {
	ej := NewEjercicioHandler(ejercicios)
	al := NewAlumnoHandler(alumnos)
	pr := NewProfesorHandler(profesores)
	ru := NewRutinaHandler(rutinas)
	re := NewRutinaEjercicioHandler(rutinas_ejercicio)
	cl := NewClaseHandler(clases)
	ch := NewClaseHorarioHandler(clases_horario)
	ins := NewInscripcionHandler(inscripciones)
	pages := web.NewHandler(ejercicios, alumnos, profesores, rutinas, rutinas_ejercicio, clases, auth)

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

	mux.HandleFunc("GET /login", pages.LoginPage)
	mux.HandleFunc("POST /login", pages.Login)
	mux.HandleFunc("POST /logout", pages.Logout)
	mux.HandleFunc("GET /", pages.Home)
	// Ejercicios, rutinas y clases: solo profesores. Profesores: solo el
	// administrador. Alumnos: profesores y administrador.
	profesor := func(next http.HandlerFunc) http.HandlerFunc {
		return pages.RequireRole(next, service.RolProfesor)
	}
	admin := func(next http.HandlerFunc) http.HandlerFunc {
		return pages.RequireRole(next, service.RolAdmin)
	}
	gestion := func(next http.HandlerFunc) http.HandlerFunc {
		return pages.RequireRole(next, service.RolProfesor, service.RolAdmin)
	}
	// La vista del alumno: solo su propia rutina.
	mux.HandleFunc("GET /mi-rutina", pages.RequireRole(pages.MiRutina, service.RolAlumno))

	mux.HandleFunc("GET /alumnos", gestion(pages.Alumnos))
	mux.HandleFunc("GET /profesores", admin(pages.Profesores))
	mux.HandleFunc("POST /alumnos", gestion(pages.CreateAlumno))
	mux.HandleFunc("POST /profesores", admin(pages.CreateProfesor))
	mux.HandleFunc("GET /alumnos/{id}/editar", gestion(pages.EditAlumno))
	mux.HandleFunc("POST /alumnos/{id}", gestion(pages.UpdateAlumno))
	mux.HandleFunc("DELETE /alumnos/{id}", gestion(pages.DeleteAlumno))
	mux.HandleFunc("GET /profesores/{id}/editar", admin(pages.EditProfesor))
	mux.HandleFunc("POST /profesores/{id}", admin(pages.UpdateProfesor))
	mux.HandleFunc("DELETE /profesores/{id}", admin(pages.DeleteProfesor))
	mux.HandleFunc("GET /ejercicios", profesor(pages.Ejercicios))
	mux.HandleFunc("POST /ejercicios", profesor(pages.CreateEjercicio))
	mux.HandleFunc("GET /ejercicios/{id}/editar", profesor(pages.EditEjercicio))
	mux.HandleFunc("POST /ejercicios/{id}", profesor(pages.UpdateEjercicio))
	mux.HandleFunc("DELETE /ejercicios/{id}", profesor(pages.DeleteEjercicio))
	mux.HandleFunc("GET /rutinas", profesor(pages.Rutinas))
	mux.HandleFunc("POST /rutinas", profesor(pages.CreateRutina))
	mux.HandleFunc("GET /rutinas/{id}/editar", profesor(pages.EditRutina))
	mux.HandleFunc("POST /rutinas/{id}", profesor(pages.UpdateRutina))
	mux.HandleFunc("DELETE /rutinas/{id}", profesor(pages.DeleteRutina))
	mux.HandleFunc("GET /clases", profesor(pages.Clases))
	mux.HandleFunc("POST /clases", profesor(pages.CreateClase))
	mux.HandleFunc("GET /clases/{id}/editar", profesor(pages.EditClase))
	mux.HandleFunc("POST /clases/{id}", profesor(pages.UpdateClase))
	mux.HandleFunc("DELETE /clases/{id}", profesor(pages.DeleteClase))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	return mux
}
