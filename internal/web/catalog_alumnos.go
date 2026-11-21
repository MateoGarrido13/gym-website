package web

import (
	"fmt"
	"net/http"
	"strconv"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

type alumnosResource struct{ h *Handler }

func (res alumnosResource) Title() string { return "Alumnos" }

func (res alumnosResource) Headers() []string {
	return []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "PLAN / VENCIMIENTO"}
}

func (res alumnosResource) Description() string {
	return "Registra alumnos y consulta sus datos principales."
}

func (res alumnosResource) Entities(r *http.Request) ([]views.Entity, error) {
	items, err := res.h.alumnos.List(r.Context())
	if err != nil {
		return nil, err
	}
	entities := make([]views.Entity, 0, len(items))
	for _, item := range items {
		detail := stringValue(item.TipoPlan)
		if item.FechaVto != nil {
			if detail != "" {
				detail += " · "
			}
			detail += "Vence " + *item.FechaVto
		}
		entities = append(entities, views.Entity{
			Values:    []string{item.Nombre + " " + item.Apellido, item.Email, stringValue(item.Telefono), detail},
			EditURL:   fmt.Sprintf("/alumnos/%d/editar", item.ID),
			DeleteURL: fmt.Sprintf("/alumnos/%d", item.ID),
		})
	}
	return entities, nil
}

func (res alumnosResource) FormOptions(r *http.Request, page *views.PageData, values views.FormValues, editing bool) error {
	page.FormValues.FechaInscripcion = defaultEnrollmentDate(values.FechaInscripcion)
	if !editing {
		return nil
	}
	routines, err := res.h.routineOptions(r)
	if err != nil {
		return err
	}
	page.Options = routines
	return nil
}

func (res alumnosResource) EditValues(r *http.Request, id int32) (views.FormValues, error) {
	item, err := res.h.alumnos.Get(r.Context(), id)
	if err != nil {
		return views.FormValues{}, err
	}
	values := views.FormValues{
		TipoPlan: stringValue(item.TipoPlan),
		FechaVto: stringValue(item.FechaVto),
	}
	if item.RutinaID != nil {
		values.RutinaID = strconv.FormatInt(int64(*item.RutinaID), 10)
	}
	return values, nil
}

func (res alumnosResource) Save(r *http.Request, id int32, update bool, values views.FormValues) error {
	if update {
		routineID, err := optionalPositiveInt32(values.RutinaID, "Selecciona una rutina válida.")
		if err != nil {
			return &formError{err.Error(), http.StatusBadRequest}
		}
		_, err = res.h.alumnos.Update(r.Context(), id, service.AlumnoUpdateInput{
			FechaVto: optional(values.FechaVto),
			TipoPlan: optional(values.TipoPlan),
			RutinaID: routineID,
		})
		return err
	}
	if err := validateIdentity(values); err != nil {
		return &formError{err.Error(), http.StatusBadRequest}
	}
	_, err := res.h.alumnos.Create(r.Context(), service.AlumnoCreateInput{
		Email:            values.Email,
		Contrasena:       values.Contrasena,
		Nombre:           values.Nombre,
		Apellido:         values.Apellido,
		Telefono:         optional(values.Telefono),
		FechaInscripcion: values.FechaInscripcion,
		FechaVto:         optional(values.FechaVto),
		TipoPlan:         optional(values.TipoPlan),
	})
	return err
}

func (res alumnosResource) Delete(r *http.Request, id int32) error {
	return res.h.alumnos.Delete(r.Context(), id)
}
