package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

// MiRutina muestra la rutina del alumno con sesión. El alumno sale siempre de
// la sesión (nunca de la URL), así nadie puede mirar la rutina de otro.
func (h *Handler) MiRutina(w http.ResponseWriter, r *http.Request) {
	user, ok := UserFrom(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	data := views.RutinaAlumnoData{UserName: user.Nombre + " " + user.Apellido}

	rutina, err := h.rutinas.DeAlumno(r.Context(), user.ID)
	switch {
	case errors.Is(err, service.ErrNotFound):
		// Sin rutina asignada: la vista muestra el mensaje correspondiente.
	case err != nil:
		log.Printf("no se pudo cargar la rutina del alumno %d: %v", user.ID, err)
		http.Error(w, "No se pudo cargar tu rutina.", http.StatusInternalServerError)
		return
	default:
		data.Nombre = rutina.Nombre
		data.Duracion = plural(int(rutina.DuracionSemanas), "semana", "semanas")
		data.Profesor = rutina.Profesor
		for _, e := range rutina.Ejercicios {
			data.Ejercicios = append(data.Ejercicios, views.EjercicioItem{
				Nombre:      e.Nombre,
				Grupo:       stringValue(e.GrupoMuscular),
				Descripcion: stringValue(e.Descripcion),
				Esquema:     esquema(e.Series, e.Repeticiones),
				Descanso:    descanso(e.DescansoSegundos),
			})
		}
	}
	h.renderAlumnoRutina(w, r, data)
}

func plural(n int, singular, pluralForm string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, pluralForm)
}

// esquema arma "4 series × 12 repeticiones" con lo que haya cargado.
func esquema(series, repeticiones *int32) string {
	var parts []string
	if series != nil {
		parts = append(parts, plural(int(*series), "serie", "series"))
	}
	if repeticiones != nil {
		parts = append(parts, plural(int(*repeticiones), "repetición", "repeticiones"))
	}
	return strings.Join(parts, " × ")
}

func descanso(segundos *int32) string {
	if segundos == nil {
		return ""
	}
	return fmt.Sprintf("%d s de descanso", *segundos)
}
