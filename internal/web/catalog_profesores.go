package web

import (
	"fmt"
	"net/http"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

type profesoresResource struct{ h *Handler }

func (res profesoresResource) Title() string { return "Profesores" }

func (res profesoresResource) Headers() []string {
	return []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "ESPECIALIDAD"}
}

func (res profesoresResource) Description() string {
	return "Registra compañeros de trabajo y consulta sus especialidades."
}

func (res profesoresResource) Entities(r *http.Request) ([]views.Entity, error) {
	items, err := res.h.profesores.List(r.Context())
	if err != nil {
		return nil, err
	}
	entities := make([]views.Entity, 0, len(items))
	for _, item := range items {
		entities = append(entities, views.Entity{
			Values:    []string{item.Nombre + " " + item.Apellido, item.Email, stringValue(item.Telefono), stringValue(item.Especialidad)},
			EditURL:   fmt.Sprintf("/profesores/%d/editar", item.ID),
			DeleteURL: fmt.Sprintf("/profesores/%d", item.ID),
		})
	}
	return entities, nil
}

func (res profesoresResource) FormOptions(r *http.Request, page *views.PageData, values views.FormValues, editing bool) error {
	return nil
}

func (res profesoresResource) EditValues(r *http.Request, id int32) (views.FormValues, error) {
	item, err := res.h.profesores.Get(r.Context(), id)
	if err != nil {
		return views.FormValues{}, err
	}
	return views.FormValues{Especialidad: stringValue(item.Especialidad)}, nil
}

func (res profesoresResource) Create(r *http.Request, values views.FormValues) error {
	if err := validateIdentity(values); err != nil {
		return &formError{err.Error(), http.StatusBadRequest}
	}
	_, err := res.h.profesores.Create(r.Context(), service.ProfesorCreateInput{
		Email:        values.Email,
		Contrasena:   values.Contrasena,
		Nombre:       values.Nombre,
		Apellido:     values.Apellido,
		Telefono:     optional(values.Telefono),
		Especialidad: optional(values.Especialidad),
	})
	return err
}

func (res profesoresResource) Update(r *http.Request, id int32, values views.FormValues) error {
	_, err := res.h.profesores.Update(r.Context(), id, service.ProfesorUpdateInput{
		Especialidad: optional(values.Especialidad),
	})
	return err
}

func (res profesoresResource) Delete(r *http.Request, id int32) error {
	return res.h.profesores.Delete(r.Context(), id)
}
