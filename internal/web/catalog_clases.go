package web

import (
	"fmt"
	"net/http"
	"strconv"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

type clasesResource struct{ h *Handler }

func (res clasesResource) Title() string { return "Clases" }

func (res clasesResource) Headers() []string {
	return []string{"CLASE", "CUPO MÁXIMO", "PROFESOR"}
}

func (res clasesResource) Description() string {
	return "Registra clases, su capacidad y el profesor responsable."
}

func (res clasesResource) Entities(r *http.Request) ([]views.Entity, error) {
	items, err := res.h.clases.List(r.Context())
	if err != nil {
		return nil, err
	}
	professors, err := res.h.professorOptions(r)
	if err != nil {
		return nil, err
	}
	entities := make([]views.Entity, 0, len(items))
	for _, item := range items {
		entities = append(entities, views.Entity{
			Values:    []string{item.Nombre, strconv.FormatInt(int64(item.MaxAlumnos), 10), professors.labels[item.ProfesorID]},
			EditURL:   fmt.Sprintf("/clases/%d/editar", item.ID),
			DeleteURL: fmt.Sprintf("/clases/%d", item.ID),
		})
	}
	return entities, nil
}

func (res clasesResource) FormOptions(r *http.Request, page *views.PageData, values views.FormValues, editing bool) error {
	professors, err := res.h.professorOptions(r)
	if err != nil {
		return err
	}
	page.Options = professors.options
	return nil
}

func (res clasesResource) EditValues(r *http.Request, id int32) (views.FormValues, error) {
	item, err := res.h.clases.Get(r.Context(), id)
	if err != nil {
		return views.FormValues{}, err
	}
	return views.FormValues{
		Nombre:     item.Nombre,
		MaxAlumnos: strconv.FormatInt(int64(item.MaxAlumnos), 10),
		ProfesorID: strconv.FormatInt(int64(item.ProfesorID), 10),
	}, nil
}

// claseFormValues valida y convierte los campos numéricos del formulario.
func claseFormValues(values views.FormValues) (maxStudents, professorID int32, err error) {
	maxStudents, err = positiveInt32(values.MaxAlumnos, "El cupo debe ser un número mayor que cero.")
	if err != nil {
		return 0, 0, &formError{err.Error(), http.StatusBadRequest}
	}
	professorID, err = positiveInt32(values.ProfesorID, "Selecciona un profesor.")
	if err != nil {
		return 0, 0, &formError{err.Error(), http.StatusBadRequest}
	}
	return maxStudents, professorID, nil
}

func (res clasesResource) Create(r *http.Request, values views.FormValues) error {
	maxStudents, professorID, err := claseFormValues(values)
	if err != nil {
		return err
	}
	_, err = res.h.clases.Create(r.Context(), service.ClaseCreateInput{
		Nombre: values.Nombre, MaxAlumnos: maxStudents, ProfesorID: professorID,
	})
	return err
}

func (res clasesResource) Update(r *http.Request, id int32, values views.FormValues) error {
	maxStudents, professorID, err := claseFormValues(values)
	if err != nil {
		return err
	}
	_, err = res.h.clases.Update(r.Context(), id, service.ClaseUpdateInput{
		Nombre: values.Nombre, MaxAlumnos: maxStudents, ProfesorID: professorID,
	})
	return err
}

func (res clasesResource) Delete(r *http.Request, id int32) error {
	return res.h.clases.Delete(r.Context(), id)
}
