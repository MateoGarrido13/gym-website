package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
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

func (h *Handler) Alumnos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	alumnos, err := h.alumnos.List(r.Context())
	if err != nil {
		log.Printf("no se pudo listar alumnos para la página: %v", err)
		http.Error(w, "No se pudieron cargar los alumnos.", http.StatusInternalServerError)
		return
	}

	entities := make([]views.Entity, 0, len(alumnos))
	for _, alumno := range alumnos {
		detalle := stringValue(alumno.TipoPlan)
		if alumno.FechaVto != nil {
			if detalle != "" {
				detalle += " · "
			}
			detalle += "Vence " + *alumno.FechaVto
		}
		entities = append(entities, views.Entity{
			ID:        alumno.ID,
			Values:    []string{alumno.Nombre + " " + alumno.Apellido, alumno.Email, stringValue(alumno.Telefono), detalle},
			EditURL:   fmt.Sprintf("/alumnos/%d/editar", alumno.ID),
			DeleteURL: fmt.Sprintf("/alumnos/%d/eliminar", alumno.ID),
		})
	}

	h.render(w, r, views.PageData{
		Title:       "Alumnos",
		Active:      "alumnos",
		Kind:        "alumnos",
		Description: "Registra alumnos y consulta sus datos principales.",
		Action:      "/alumnos",
		Headers:     []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "PLAN / VENCIMIENTO"},
		ShowActions: true,
		FormValues:  views.FormValues{FechaInscripcion: time.Now().Format(time.DateOnly)},
		Entities:    entities,
	})
}

func (h *Handler) Profesores(w http.ResponseWriter, r *http.Request) {
	profesores, err := h.profesores.List(r.Context())
	if err != nil {
		log.Printf("no se pudo listar profesores para la página: %v", err)
		http.Error(w, "No se pudieron cargar los profesores.", http.StatusInternalServerError)
		return
	}

	entities := make([]views.Entity, 0, len(profesores))
	for _, profesor := range profesores {
		entities = append(entities, views.Entity{
			ID:        profesor.ID,
			Values:    []string{profesor.Nombre + " " + profesor.Apellido, profesor.Email, stringValue(profesor.Telefono), stringValue(profesor.Especialidad)},
			EditURL:   fmt.Sprintf("/profesores/%d/editar", profesor.ID),
			DeleteURL: fmt.Sprintf("/profesores/%d/eliminar", profesor.ID),
		})
	}

	h.render(w, r, views.PageData{
		Title:       "Profesores",
		Active:      "profesores",
		Kind:        "profesores",
		Description: "Registra compañeros de trabajo y consulta sus especialidades.",
		Action:      "/profesores",
		Headers:     []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "ESPECIALIDAD"},
		ShowActions: true,
		Entities:    entities,
	})
}

func (h *Handler) CreateAlumno(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario.", http.StatusBadRequest)
		return
	}
	values := formValues(r)
	if err := validateIdentity(values); err != nil {
		h.renderFormError(w, r, "alumnos", values, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.alumnos.Create(r.Context(), service.AlumnoCreateInput{
		Email:            values.Email,
		Contrasena:       values.Contrasena,
		Nombre:           values.Nombre,
		Apellido:         values.Apellido,
		Telefono:         optional(values.Telefono),
		FechaInscripcion: values.FechaInscripcion,
		FechaVto:         optional(values.FechaVto),
		TipoPlan:         optional(values.TipoPlan),
	})
	if err != nil {
		message, status := formError(err)
		h.renderFormError(w, r, "alumnos", values, message, status)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) CreateProfesor(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario.", http.StatusBadRequest)
		return
	}
	values := formValues(r)
	if err := validateIdentity(values); err != nil {
		h.renderFormError(w, r, "profesores", values, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.profesores.Create(r.Context(), service.ProfesorCreateInput{
		Email:        values.Email,
		Contrasena:   values.Contrasena,
		Nombre:       values.Nombre,
		Apellido:     values.Apellido,
		Telefono:     optional(values.Telefono),
		Especialidad: optional(values.Especialidad),
	})
	if err != nil {
		message, status := formError(err)
		h.renderFormError(w, r, "profesores", values, message, status)
		return
	}
	http.Redirect(w, r, "/profesores", http.StatusSeeOther)
}
