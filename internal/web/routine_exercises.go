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
	localIDs := make(map[int32]bool, len(available))
	for _, exercise := range available {
		localIDs[exercise.ID] = true
	}
	selected := make([]int32, 0, len(values))
	seen := make(map[int32]bool, len(values))
	for _, value := range values {
		id, err := positiveInt32(value, "La selección de ejercicios no es válida.")
		if err != nil {
			return nil, errors.New("La selección de ejercicios no es válida. Actualiza la página e inténtalo nuevamente.")
		}
		if seen[id] {
			return nil, errors.New("La selección de ejercicios no es válida. Actualiza la página e inténtalo nuevamente.")
		}
		seen[id] = true

		// 1) Si el ID ya existe localmente, se usa directamente.
		if localIDs[id] {
			selected = append(selected, id)
			continue
		}

		// 2) Si no existe localmente, consultar wger para validar y sincronizar.
		if h.wgerACL == nil {
			return nil, errors.New("La selección de ejercicios no es válida. Actualiza la página e inténtalo nuevamente.")
		}
		wex, err := h.wgerACL.GetExerciseByID(r.Context(), int(id))
		if err != nil || wex == nil {
			return nil, errors.New("La selección de ejercicios no es válida. Actualiza la página e inténtalo nuevamente.")
		}

		// 3) Sincronizar el ejercicio de wger a la tabla local para respetar la FK.
		localEx, err := h.ejercicios.FindOrCreateByNombreGrupo(
			r.Context(), wex.Name, wex.MuscleGroup, &wex.Description,
		)
		if err != nil {
			return nil, errors.New("No se pudo sincronizar el ejercicio seleccionado. Inténtalo nuevamente.")
		}

		if seen[localEx.ID] {
			return nil, errors.New("La selección de ejercicios no es válida. Actualiza la página e inténtalo nuevamente.")
		}
		selected = append(selected, localEx.ID)
	}
	if len(selected) == 0 {
		return nil, errors.New("No se seleccionaron ejercicios. Por favor, selecciona al menos un ejercicio.")
	}
	return selected, nil
}

type routineExercise struct {
	ID           int32
	Series       *int32
	Repeticiones *int32
	Descanso     *int32
}

func (h *Handler) syncRoutineExercises(r *http.Request, routineID int32, selected []routineExercise) error {
	current, err := h.rutinaEjercicios.List(r.Context(), routineID)
	if err != nil {
		return err
	}
	existing := make(map[int32]bool, len(current))
	for _, item := range current {
		existing[item.EjercicioID] = true
	}
	wanted := make(map[int32]bool, len(selected))
	for _, ex := range selected {
		wanted[ex.ID] = true
	}
	for exerciseID := range existing {
		if !wanted[exerciseID] {
			if err := h.rutinaEjercicios.Delete(r.Context(), routineID, exerciseID); err != nil {
				return err
			}
		}
	}
	for i, ex := range selected {
		orden := int32(i + 1)
		if existing[ex.ID] {
			if _, err := h.rutinaEjercicios.Update(r.Context(), routineID, ex.ID, service.RutinaEjercicioUpdateInput{
				Orden: orden, Series: ex.Series, Repeticiones: ex.Repeticiones, DescansoSegundos: ex.Descanso,
			}); err != nil {
				return err
			}
			continue
		}
		if _, err := h.rutinaEjercicios.Create(r.Context(), routineID, service.RutinaEjercicioCreateInput{
			EjercicioID: ex.ID, Orden: orden, Series: ex.Series, Repeticiones: ex.Repeticiones, DescansoSegundos: ex.Descanso,
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
