package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

type rutinasResource struct{ h *Handler }

func (res rutinasResource) Title() string { return "Rutinas" }

func (res rutinasResource) Headers() []string {
	return []string{"RUTINA", "DURACIÓN", "PROFESOR"}
}

func (res rutinasResource) Description() string {
	return "Define planes de entrenamiento y asígnalos a un profesor."
}

func (res rutinasResource) Entities(r *http.Request) ([]views.Entity, error) {
	items, err := res.h.rutinas.List(r.Context())
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
			Values:    []string{item.Nombre, fmt.Sprintf("%d semanas", item.DuracionSemanas), professors.labels[item.ProfesorID]},
			EditURL:   fmt.Sprintf("/rutinas/%d/editar", item.ID),
			DeleteURL: fmt.Sprintf("/rutinas/%d", item.ID),
		})
	}
	return entities, nil
}

func (res rutinasResource) FormOptions(r *http.Request, page *views.PageData, values views.FormValues, editing bool) error {
	professors, err := res.h.professorOptions(r)
	if err != nil {
		return err
	}
	page.Options = professors.options
	exercises, err := res.h.exerciseOptions(r)
	if err != nil {
		return err
	}
	page.ExerciseOptions = exercises
	page.SelectedExercises = values.SelectedExercises
	return nil
}

func (res rutinasResource) EditValues(r *http.Request, id int32) (views.FormValues, error) {
	item, err := res.h.rutinas.Get(r.Context(), id)
	if err != nil {
		return views.FormValues{}, err
	}
	values := views.FormValues{
		Nombre:          item.Nombre,
		DuracionSemanas: strconv.FormatInt(int64(item.DuracionSemanas), 10),
		ProfesorID:      strconv.FormatInt(int64(item.ProfesorID), 10),
	}
	exercises, err := res.h.rutinaEjercicios.List(r.Context(), id)
	if err != nil {
		return views.FormValues{}, err
	}
	for _, exercise := range exercises {
		values.SelectedExercises = append(values.SelectedExercises, strconv.FormatInt(int64(exercise.EjercicioID), 10))
	}
	return values, nil
}

// formValues valida los campos del formulario de rutinas.
func (res rutinasResource) formValues(r *http.Request, values views.FormValues) (weeks, professorID int32, exerciseIDs []int32, err error) {
	weeks, err = positiveInt32(values.DuracionSemanas, "La duración debe ser un número mayor que cero.")
	if err != nil {
		return 0, 0, nil, &formError{err.Error(), http.StatusBadRequest}
	}
	professorID, err = positiveInt32(values.ProfesorID, "Selecciona un profesor.")
	if err != nil {
		return 0, 0, nil, &formError{err.Error(), http.StatusBadRequest}
	}
	exerciseIDs, err = res.h.validateRoutineExerciseIDs(r, values.SelectedExercises)
	if err != nil {
		return 0, 0, nil, &formError{err.Error(), http.StatusBadRequest}
	}
	return weeks, professorID, exerciseIDs, nil
}

func (res rutinasResource) Create(r *http.Request, values views.FormValues) error {
	weeks, professorID, exerciseIDs, err := res.formValues(r, values)
	if err != nil {
		return err
	}
	created, err := res.h.rutinas.Create(r.Context(), service.RutinaCreateInput{
		Nombre: values.Nombre, DuracionSemanas: weeks, ProfesorID: professorID,
	})
	if err != nil {
		return err
	}
	if err := res.h.syncRoutineExercises(r, created.ID, exerciseIDs); err != nil {
		if cleanupErr := res.h.rutinas.Delete(r.Context(), created.ID); cleanupErr != nil {
			log.Printf("no se pudo limpiar la rutina %d tras fallar la asociación de ejercicios: %v", created.ID, cleanupErr)
		}
		return err
	}
	return nil
}

func (res rutinasResource) Update(r *http.Request, id int32, values views.FormValues) error {
	weeks, professorID, exerciseIDs, err := res.formValues(r, values)
	if err != nil {
		return err
	}
	if _, err := res.h.rutinas.Update(r.Context(), id, service.RutinaUpdateInput{
		Nombre: values.Nombre, DuracionSemanas: weeks, ProfesorID: professorID,
	}); err != nil {
		return err
	}
	return res.h.syncRoutineExercises(r, id, exerciseIDs)
}

func (res rutinasResource) Delete(r *http.Request, id int32) error {
	return res.h.rutinas.Delete(r.Context(), id)
}
