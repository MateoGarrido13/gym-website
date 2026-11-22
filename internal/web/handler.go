package web

import (
	"os"

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
	secureCookies    bool
}

// NewHandler arma el handler de páginas. La cookie de sesión lleva el flag
// Secure solo si COOKIE_SECURE=true: en local no hay HTTPS y el navegador
// descartaría la cookie. En producción, detrás de HTTPS, debe activarse.
func NewHandler(ejercicios *service.EjercicioService, alumnos *service.AlumnoService, profesores *service.ProfesorService, rutinas *service.RutinaService, rutinaEjercicios *service.RutinaEjercicioService, clases *service.ClaseService, auth *service.AuthService) *Handler {
	return &Handler{
		ejercicios: ejercicios, alumnos: alumnos, profesores: profesores,
		rutinas: rutinas, rutinaEjercicios: rutinaEjercicios, clases: clases,
		auth: auth, secureCookies: os.Getenv("COOKIE_SECURE") == "true",
	}
}
