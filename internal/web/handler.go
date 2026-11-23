package web

import (
	"net/http"
	"os"

	"PRACTICO_DOS/internal/integration"
	"PRACTICO_DOS/internal/service"
)

type Handler struct {
	ejercicios       *service.EjercicioService
	alumnos          *service.AlumnoService
	profesores       *service.ProfesorService
	rutinas          *service.RutinaService
	rutinaEjercicios *service.RutinaEjercicioService
	clases           *service.ClaseService
	auth             *service.AuthService
	wgerACL          integration.WgerACL
	secureCookies    bool
}

// NewHandler arma el handler de páginas. La cookie de sesión lleva el flag
// Secure solo si COOKIE_SECURE=true: en local no hay HTTPS y el navegador
// descartaría la cookie. En producción, detrás de HTTPS, debe activarse.
func NewHandler(ejercicios *service.EjercicioService, alumnos *service.AlumnoService, profesores *service.ProfesorService, rutinas *service.RutinaService, rutinaEjercicios *service.RutinaEjercicioService, clases *service.ClaseService, auth *service.AuthService, wgerACL integration.WgerACL) *Handler {
	return &Handler{
		ejercicios: ejercicios, alumnos: alumnos, profesores: profesores,
		rutinas: rutinas, rutinaEjercicios: rutinaEjercicios, clases: clases,
		auth: auth, wgerACL: wgerACL, secureCookies: os.Getenv("COOKIE_SECURE") == "true",
	}
}

// SearchRoutineExercises expone el handler HTMX de búsqueda de ejercicios de wger
// para ser registrado como ruta desde el router.
func (h *Handler) SearchRoutineExercises(w http.ResponseWriter, r *http.Request) {
	SearchRoutineExercisesHandler(h.wgerACL)(w, r)
}
