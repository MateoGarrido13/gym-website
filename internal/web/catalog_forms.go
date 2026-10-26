package web

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

func (h *Handler) editCatalog(w http.ResponseWriter, r *http.Request, kind string) {
	id, err := catalogID(r)
	if err != nil {
		http.Error(w, "El identificador no es válido.", http.StatusBadRequest)
		return
	}

	values := views.FormValues{}
	switch kind {
	case "alumnos":
		item, err := h.alumnos.Get(r.Context(), id)
		if err != nil {
			h.catalogRequestError(w, r, kind, err)
			return
		}
		values.TipoPlan = stringValue(item.TipoPlan)
		values.FechaVto = stringValue(item.FechaVto)
		if item.RutinaID != nil {
			values.RutinaID = strconv.FormatInt(int64(*item.RutinaID), 10)
		}
	case "profesores":
		item, err := h.profesores.Get(r.Context(), id)
		if err != nil {
			h.catalogRequestError(w, r, kind, err)
			return
		}
		values.Especialidad = stringValue(item.Especialidad)
	case "ejercicios":
		item, err := h.ejercicios.Get(r.Context(), id)
		if err != nil {
			h.catalogRequestError(w, r, kind, err)
			return
		}
		values.Nombre = item.Nombre
		values.Descripcion = stringValue(item.Descripcion)
		values.GrupoMuscular = stringValue(item.GrupoMuscular)
	case "rutinas":
		item, err := h.rutinas.Get(r.Context(), id)
		if err != nil {
			h.catalogRequestError(w, r, kind, err)
			return
		}
		values.Nombre = item.Nombre
		values.DuracionSemanas = strconv.FormatInt(int64(item.DuracionSemanas), 10)
		values.ProfesorID = strconv.FormatInt(int64(item.ProfesorID), 10)
		exercises, err := h.rutinaEjercicios.List(r.Context(), id)
		if err != nil {
			log.Printf("no se pudieron cargar los ejercicios de la rutina %d: %v", id, err)
			http.Error(w, "No se pudieron cargar los ejercicios de la rutina.", http.StatusInternalServerError)
			return
		}
		for _, exercise := range exercises {
			values.SelectedExercises = append(values.SelectedExercises, strconv.FormatInt(int64(exercise.EjercicioID), 10))
		}
	case "clases":
		item, err := h.clases.Get(r.Context(), id)
		if err != nil {
			h.catalogRequestError(w, r, kind, err)
			return
		}
		values.Nombre = item.Nombre
		values.MaxAlumnos = strconv.FormatInt(int64(item.MaxAlumnos), 10)
		values.ProfesorID = strconv.FormatInt(int64(item.ProfesorID), 10)
	default:
		http.NotFound(w, r)
		return
	}
	h.showCatalog(w, r, kind, values, true, id, "", http.StatusOK)
}

func (h *Handler) saveCatalog(w http.ResponseWriter, r *http.Request, kind string, update bool) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario.", http.StatusBadRequest)
		return
	}
	values := views.FormValues{
		Nombre:            strings.TrimSpace(r.FormValue("nombre")),
		FechaVto:          strings.TrimSpace(r.FormValue("fecha_vto")),
		TipoPlan:          strings.TrimSpace(r.FormValue("tipo_plan")),
		Especialidad:      strings.TrimSpace(r.FormValue("especialidad")),
		Descripcion:       strings.TrimSpace(r.FormValue("descripcion")),
		GrupoMuscular:     strings.TrimSpace(r.FormValue("grupo_muscular")),
		DuracionSemanas:   strings.TrimSpace(r.FormValue("duracion_semanas")),
		ProfesorID:        strings.TrimSpace(r.FormValue("profesor_id")),
		MaxAlumnos:        strings.TrimSpace(r.FormValue("max_alumnos")),
		RutinaID:          strings.TrimSpace(r.FormValue("rutina_id")),
		SelectedExercises: append([]string(nil), r.Form["ejercicio_id"]...),
	}

	id := int32(0)
	if update {
		var err error
		id, err = catalogID(r)
		if err != nil {
			http.Error(w, "El identificador no es válido.", http.StatusBadRequest)
			return
		}
	}

	var err error
	switch kind {
	case "alumnos":
		routineID, parseErr := optionalPositiveInt32(values.RutinaID, "Selecciona una rutina válida.")
		if parseErr != nil {
			h.showCatalog(w, r, kind, values, true, id, parseErr.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.alumnos.Update(r.Context(), id, service.AlumnoUpdateInput{
			FechaVto: optional(values.FechaVto),
			TipoPlan: optional(values.TipoPlan),
			RutinaID: routineID,
		})
	case "profesores":
		_, err = h.profesores.Update(r.Context(), id, service.ProfesorUpdateInput{
			Especialidad: optional(values.Especialidad),
		})
	case "ejercicios":
		input := service.EjercicioInput{
			Nombre:        values.Nombre,
			Descripcion:   optional(values.Descripcion),
			GrupoMuscular: optional(values.GrupoMuscular),
		}
		if update {
			_, err = h.ejercicios.Update(r.Context(), id, input)
		} else {
			_, err = h.ejercicios.Create(r.Context(), input)
		}
	case "rutinas":
		weeks, parseErr := positiveInt32(values.DuracionSemanas, "La duración debe ser un número mayor que cero.")
		if parseErr != nil {
			h.showCatalog(w, r, kind, values, update, id, parseErr.Error(), http.StatusBadRequest)
			return
		}
		professorID, parseErr := positiveInt32(values.ProfesorID, "Selecciona un profesor.")
		if parseErr != nil {
			h.showCatalog(w, r, kind, values, update, id, parseErr.Error(), http.StatusBadRequest)
			return
		}
		exerciseIDs, parseErr := h.validateRoutineExerciseIDs(r, values.SelectedExercises)
		if parseErr != nil {
			h.showCatalog(w, r, kind, values, update, id, parseErr.Error(), http.StatusBadRequest)
			return
		}
		if update {
			_, err = h.rutinas.Update(r.Context(), id, service.RutinaUpdateInput{
				Nombre: values.Nombre, DuracionSemanas: weeks, ProfesorID: professorID,
			})
			if err == nil {
				err = h.syncRoutineExercises(r, id, exerciseIDs)
			}
		} else {
			var created service.Rutina
			created, err = h.rutinas.Create(r.Context(), service.RutinaCreateInput{
				Nombre: values.Nombre, DuracionSemanas: weeks, ProfesorID: professorID,
			})
			if err == nil {
				err = h.syncRoutineExercises(r, created.ID, exerciseIDs)
				if err != nil {
					if cleanupErr := h.rutinas.Delete(r.Context(), created.ID); cleanupErr != nil {
						log.Printf("no se pudo limpiar la rutina %d tras fallar la asociación de ejercicios: %v", created.ID, cleanupErr)
					}
				}
			}
		}
	case "clases":
		maxStudents, parseErr := positiveInt32(values.MaxAlumnos, "El cupo debe ser un número mayor que cero.")
		if parseErr != nil {
			h.showCatalog(w, r, kind, values, update, id, parseErr.Error(), http.StatusBadRequest)
			return
		}
		professorID, parseErr := positiveInt32(values.ProfesorID, "Selecciona un profesor.")
		if parseErr != nil {
			h.showCatalog(w, r, kind, values, update, id, parseErr.Error(), http.StatusBadRequest)
			return
		}
		if update {
			_, err = h.clases.Update(r.Context(), id, service.ClaseUpdateInput{
				Nombre: values.Nombre, MaxAlumnos: maxStudents, ProfesorID: professorID,
			})
		} else {
			_, err = h.clases.Create(r.Context(), service.ClaseCreateInput{
				Nombre: values.Nombre, MaxAlumnos: maxStudents, ProfesorID: professorID,
			})
		}
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		message, status := catalogFormError(err)
		h.showCatalog(w, r, kind, values, update, id, message, status)
		return
	}
	http.Redirect(w, r, catalogPathForRedirect(kind), http.StatusSeeOther)
}

func (h *Handler) deleteCatalog(w http.ResponseWriter, r *http.Request, kind string) {
	id, err := catalogID(r)
	if err != nil {
		http.Error(w, "El identificador no es válido.", http.StatusBadRequest)
		return
	}

	switch kind {
	case "alumnos":
		err = h.alumnos.Delete(r.Context(), id)
	case "profesores":
		err = h.profesores.Delete(r.Context(), id)
	case "ejercicios":
		err = h.ejercicios.Delete(r.Context(), id)
	case "rutinas":
		err = h.rutinas.Delete(r.Context(), id)
	case "clases":
		err = h.clases.Delete(r.Context(), id)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		message, status := catalogFormError(err)
		h.showCatalog(w, r, kind, views.FormValues{}, false, 0, message, status)
		return
	}
	http.Redirect(w, r, catalogPathForRedirect(kind), http.StatusSeeOther)
}
