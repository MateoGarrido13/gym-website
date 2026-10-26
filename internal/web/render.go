package web

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"PRACTICO_DOS/views"
)

func (h *Handler) renderFormError(w http.ResponseWriter, r *http.Request, kind string, values submittedValues, message string, status int) {
	title, description, action := "Alumnos", "Registra alumnos y consulta sus datos principales.", "/alumnos"
	active := "alumnos"
	if kind == "profesores" {
		title, description, action = "Profesores", "Registra compañeros de trabajo y consulta sus especialidades.", "/profesores"
		active = "profesores"
	}
	page := views.PageData{
		Title:       title,
		Active:      active,
		Kind:        kind,
		Description: description,
		Action:      action,
		Headers:     []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "PLAN / VENCIMIENTO"},
		ShowActions: true,
		FormError:   message,
		FormValues:  views.FormValues{Nombre: values.Nombre, Apellido: values.Apellido, Email: values.Email, Telefono: values.Telefono, FechaInscripcion: defaultEnrollmentDate(values.FechaInscripcion), FechaVto: values.FechaVto, TipoPlan: values.TipoPlan, Especialidad: values.Especialidad},
	}
	if kind == "profesores" {
		page.Headers = []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "ESPECIALIDAD"}
	}

	var entities []views.Entity
	if kind == "alumnos" {
		alumnos, err := h.alumnos.List(r.Context())
		if err != nil {
			log.Printf("no se pudo volver a listar alumnos después de un error de formulario: %v", err)
			http.Error(w, "No se pudieron cargar los alumnos.", http.StatusInternalServerError)
			return
		}
		entities = make([]views.Entity, 0, len(alumnos))
		for _, alumno := range alumnos {
			detail := stringValue(alumno.TipoPlan)
			if alumno.FechaVto != nil {
				if detail != "" {
					detail += " · "
				}
				detail += "Vence " + *alumno.FechaVto
			}
			entities = append(entities, views.Entity{
				ID:        alumno.ID,
				Values:    []string{alumno.Nombre + " " + alumno.Apellido, alumno.Email, stringValue(alumno.Telefono), detail},
				EditURL:   fmt.Sprintf("/alumnos/%d/editar", alumno.ID),
				DeleteURL: fmt.Sprintf("/alumnos/%d/eliminar", alumno.ID),
			})
		}
	} else {
		profesores, err := h.profesores.List(r.Context())
		if err != nil {
			log.Printf("no se pudo volver a listar profesores después de un error de formulario: %v", err)
			http.Error(w, "No se pudieron cargar los profesores.", http.StatusInternalServerError)
			return
		}
		entities = make([]views.Entity, 0, len(profesores))
		for _, profesor := range profesores {
			entities = append(entities, views.Entity{
				ID:        profesor.ID,
				Values:    []string{profesor.Nombre + " " + profesor.Apellido, profesor.Email, stringValue(profesor.Telefono), stringValue(profesor.Especialidad)},
				EditURL:   fmt.Sprintf("/profesores/%d/editar", profesor.ID),
				DeleteURL: fmt.Sprintf("/profesores/%d/eliminar", profesor.ID),
			})
		}
	}
	page.Entities = entities
	h.renderStatus(w, r, page, status)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, page views.PageData) {
	h.renderStatus(w, r, page, http.StatusOK)
}

func (h *Handler) renderStatus(w http.ResponseWriter, r *http.Request, page views.PageData, status int) {
	var output bytes.Buffer
	if err := views.IndexPage(page).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar la página %q: %v", page.Title, err)
		http.Error(w, "No se pudo renderizar la página.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := output.WriteTo(w); err != nil {
		log.Printf("no se pudo escribir la página %q: %v", page.Title, err)
	}
}
