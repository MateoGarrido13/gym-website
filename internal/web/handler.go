package web

import (
	"PRACTICO_DOS/internal/service"
)

type Handler struct {
	ejercicios       *service.EjercicioService
	alumnos          *service.AlumnoService
	profesores       *service.ProfesorService
	rutinas          *service.RutinaService
	rutinaEjercicios *service.RutinaEjercicioService
	clases           *service.ClaseService
}

func NewHandler(ejercicios *service.EjercicioService, alumnos *service.AlumnoService, profesores *service.ProfesorService, rutinas *service.RutinaService, rutinaEjercicios *service.RutinaEjercicioService, clases *service.ClaseService) *Handler {
	return &Handler{ejercicios: ejercicios, alumnos: alumnos, profesores: profesores, rutinas: rutinas, rutinaEjercicios: rutinaEjercicios, clases: clases}
}
