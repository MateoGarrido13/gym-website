package web

import (
	"fmt"
	"net/http"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

type ejerciciosResource struct{ h *Handler }

func (res ejerciciosResource) Title() string { return "Ejercicios" }

func (res ejerciciosResource) Headers() []string {
	return []string{"EJERCICIO", "GRUPO MUSCULAR", "DESCRIPCIÓN"}
}

func (res ejerciciosResource) Description() string {
	return "Crea y administra el catálogo de ejercicios disponibles para las rutinas."
}

func (res ejerciciosResource) Entities(r *http.Request) ([]views.Entity, error) {
	items, err := res.h.ejercicios.List(r.Context())
	if err != nil {
		return nil, err
	}
	entities := make([]views.Entity, 0, len(items))
	for _, item := range items {
		entities = append(entities, views.Entity{
			Values:    []string{item.Nombre, stringValue(item.GrupoMuscular), stringValue(item.Descripcion)},
			EditURL:   fmt.Sprintf("/ejercicios/%d/editar", item.ID),
			DeleteURL: fmt.Sprintf("/ejercicios/%d", item.ID),
		})
	}
	return entities, nil
}

func (res ejerciciosResource) FormOptions(r *http.Request, page *views.PageData, values views.FormValues, editing bool) error {
	return nil
}

func (res ejerciciosResource) EditValues(r *http.Request, id int32) (views.FormValues, error) {
	item, err := res.h.ejercicios.Get(r.Context(), id)
	if err != nil {
		return views.FormValues{}, err
	}
	return views.FormValues{
		Nombre:        item.Nombre,
		Descripcion:   stringValue(item.Descripcion),
		GrupoMuscular: stringValue(item.GrupoMuscular),
	}, nil
}

func ejercicioInput(values views.FormValues) service.EjercicioInput {
	return service.EjercicioInput{
		Nombre:        values.Nombre,
		Descripcion:   optional(values.Descripcion),
		GrupoMuscular: optional(values.GrupoMuscular),
	}
}

func (res ejerciciosResource) Create(r *http.Request, values views.FormValues) error {
	_, err := res.h.ejercicios.Create(r.Context(), ejercicioInput(values))
	return err
}

func (res ejerciciosResource) Update(r *http.Request, id int32, values views.FormValues) error {
	_, err := res.h.ejercicios.Update(r.Context(), id, ejercicioInput(values))
	return err
}

func (res ejerciciosResource) Delete(r *http.Request, id int32) error {
	return res.h.ejercicios.Delete(r.Context(), id)
}
