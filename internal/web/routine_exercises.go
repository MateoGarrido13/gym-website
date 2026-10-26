package web

import (
	"errors"
	"fmt"
	"net/http"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

func (h *Handler) exerciseOptions(r *http.Request) ([]views.Option, error) {
	exercises, err := h.ejercicios.List(r.Context())
	if err != nil {
		return nil, err
	}
	options := make([]views.Option, 0, len(exercises))
	for _, exercise := range exercises {
		label := exercise.Nombre
		if muscle := stringValue(exercise.GrupoMuscular); muscle != "" {
			label += " · " + muscle
		}
		options = append(options, views.Option{ID: exercise.ID, Label: label})
	}
	return options, nil
}

func (h *Handler) validateRoutineExerciseIDs(r *http.Request, values []string) ([]int32, error) {
	available, err := h.exerciseOptions(r)
	if err != nil {
		return nil, err
	}
	allowed := make(map[int32]bool, len(available))
	for _, exercise := range available {
		allowed[exercise.ID] = true
	}
	selected := make([]int32, 0, len(values))
	seen := make(map[int32]bool, len(values))
	for _, value := range values {
		id, err := positiveInt32(value, "La selección de ejercicios no es válida.")
		if err != nil || !allowed[id] || seen[id] {
			return nil, errors.New("La selección de ejercicios no es válida. Actualiza la página e inténtalo nuevamente.")
		}
		selected = append(selected, id)
		seen[id] = true
	}
	return selected, nil
}

func (h *Handler) syncRoutineExercises(r *http.Request, routineID int32, selected []int32) error {
	current, err := h.rutinaEjercicios.List(r.Context(), routineID)
	if err != nil {
		return err
	}
	existing := make(map[int32]service.RutinaEjercicio, len(current))
	var maxOrder int32
	for _, item := range current {
		existing[item.EjercicioID] = item
		if item.Orden > maxOrder {
			maxOrder = item.Orden
		}
	}
	wanted := make(map[int32]bool, len(selected))
	for _, exerciseID := range selected {
		wanted[exerciseID] = true
	}
	for exerciseID := range existing {
		if !wanted[exerciseID] {
			if err := h.rutinaEjercicios.Delete(r.Context(), routineID, exerciseID); err != nil {
				return err
			}
		}
	}
	for _, exerciseID := range selected {
		if _, alreadyAdded := existing[exerciseID]; alreadyAdded {
			continue
		}
		maxOrder++
		if _, err := h.rutinaEjercicios.Create(r.Context(), routineID, service.RutinaEjercicioCreateInput{
			EjercicioID: exerciseID,
			Orden:       maxOrder,
		}); err != nil {
			return err
		}
	}
	return nil
}

type professorDirectory struct {
	options []views.Option
	labels  map[int32]string
}

func (h *Handler) professorOptions(r *http.Request) (professorDirectory, error) {
	professors, err := h.profesores.List(r.Context())
	if err != nil {
		return professorDirectory{}, err
	}
	directory := professorDirectory{
		options: make([]views.Option, 0, len(professors)),
		labels:  make(map[int32]string, len(professors)),
	}
	for _, professor := range professors {
		label := professor.Nombre + " " + professor.Apellido
		directory.options = append(directory.options, views.Option{ID: professor.ID, Label: label})
		directory.labels[professor.ID] = label
	}
	return directory, nil
}

func (h *Handler) routineOptions(r *http.Request) ([]views.Option, error) {
	routines, err := h.rutinas.List(r.Context())
	if err != nil {
		return nil, err
	}
	options := make([]views.Option, 0, len(routines))
	for _, routine := range routines {
		options = append(options, views.Option{
			ID:    routine.ID,
			Label: fmt.Sprintf("%s (%d semanas)", routine.Nombre, routine.DuracionSemanas),
		})
	}
	return options, nil
}
