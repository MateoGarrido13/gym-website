package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

func (h *Handler) showCatalog(w http.ResponseWriter, r *http.Request, kind string, values views.FormValues, editing bool, id int32, message string, status int) {
	page, err := h.catalogPage(r, kind, values, editing, id, message)
	if err != nil {
		log.Printf("no se pudo cargar %s para la página: %v", kind, err)
		http.Error(w, "No se pudieron cargar los registros.", http.StatusInternalServerError)
		return
	}
	h.renderStatus(w, r, page, status)
}

func (h *Handler) catalogPage(r *http.Request, kind string, values views.FormValues, editing bool, id int32, message string) (views.PageData, error) {
	page := views.PageData{
		Active:      kind,
		Kind:        kind,
		Action:      "/" + kind,
		Description: catalogDescription(kind),
		FormError:   message,
		IsEditing:   editing,
		ShowActions: true,
		FormValues:  values,
	}
	switch kind {
	case "alumnos":
		page.Title = "Alumnos"
		page.Headers = []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "PLAN / VENCIMIENTO"}
		items, err := h.alumnos.List(r.Context())
		if err != nil {
			return page, err
		}
		page.Entities = make([]views.Entity, 0, len(items))
		for _, item := range items {
			detail := stringValue(item.TipoPlan)
			if item.FechaVto != nil {
				if detail != "" {
					detail += " · "
				}
				detail += "Vence " + *item.FechaVto
			}
			page.Entities = append(page.Entities, views.Entity{
				Values:  []string{item.Nombre + " " + item.Apellido, item.Email, stringValue(item.Telefono), detail},
				EditURL: fmt.Sprintf("/alumnos/%d/editar", item.ID), DeleteURL: fmt.Sprintf("/alumnos/%d/eliminar", item.ID),
			})
		}
		if editing {
			routines, err := h.routineOptions(r)
			if err != nil {
				return page, err
			}
			page.Options = routines
		}
	case "profesores":
		page.Title = "Profesores"
		page.Headers = []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "ESPECIALIDAD"}
		items, err := h.profesores.List(r.Context())
		if err != nil {
			return page, err
		}
		page.Entities = make([]views.Entity, 0, len(items))
		for _, item := range items {
			page.Entities = append(page.Entities, views.Entity{
				Values:  []string{item.Nombre + " " + item.Apellido, item.Email, stringValue(item.Telefono), stringValue(item.Especialidad)},
				EditURL: fmt.Sprintf("/profesores/%d/editar", item.ID), DeleteURL: fmt.Sprintf("/profesores/%d/eliminar", item.ID),
			})
		}
	case "ejercicios":
		page.Title = "Ejercicios"
		page.Headers = []string{"EJERCICIO", "GRUPO MUSCULAR", "DESCRIPCIÓN"}
		items, err := h.ejercicios.List(r.Context())
		if err != nil {
			return page, err
		}
		page.Entities = make([]views.Entity, 0, len(items))
		for _, item := range items {
			page.Entities = append(page.Entities, views.Entity{
				Values:  []string{item.Nombre, stringValue(item.GrupoMuscular), stringValue(item.Descripcion)},
				EditURL: fmt.Sprintf("/ejercicios/%d/editar", item.ID), DeleteURL: fmt.Sprintf("/ejercicios/%d/eliminar", item.ID),
			})
		}
	case "rutinas":
		page.Title = "Rutinas"
		page.Headers = []string{"RUTINA", "DURACIÓN", "PROFESOR"}
		items, err := h.rutinas.List(r.Context())
		if err != nil {
			return page, err
		}
		professors, err := h.professorOptions(r)
		if err != nil {
			return page, err
		}
		page.Options = professors.options
		exercises, err := h.exerciseOptions(r)
		if err != nil {
			return page, err
		}
		page.ExerciseOptions = exercises
		page.SelectedExercises = values.SelectedExercises
		page.Entities = make([]views.Entity, 0, len(items))
		for _, item := range items {
			page.Entities = append(page.Entities, views.Entity{
				Values:  []string{item.Nombre, fmt.Sprintf("%d semanas", item.DuracionSemanas), professors.labels[item.ProfesorID]},
				EditURL: fmt.Sprintf("/rutinas/%d/editar", item.ID), DeleteURL: fmt.Sprintf("/rutinas/%d/eliminar", item.ID),
			})
		}

	case "clases":
		page.Title = "Clases"
		page.Headers = []string{"CLASE", "CUPO MÁXIMO", "PROFESOR"}
		items, err := h.clases.List(r.Context())
		if err != nil {
			return page, err
		}
		professors, err := h.professorOptions(r)
		if err != nil {
			return page, err
		}
		page.Options = professors.options
		page.Entities = make([]views.Entity, 0, len(items))
		for _, item := range items {
			page.Entities = append(page.Entities, views.Entity{
				Values:  []string{item.Nombre, strconv.FormatInt(int64(item.MaxAlumnos), 10), professors.labels[item.ProfesorID]},
				EditURL: fmt.Sprintf("/clases/%d/editar", item.ID), DeleteURL: fmt.Sprintf("/clases/%d/eliminar", item.ID),
			})
		}
	default:
		return page, service.ErrNotFound
	}
	if editing {
		page.Action = fmt.Sprintf("/%s/%d", kind, id)
	}
	return page, nil
}
