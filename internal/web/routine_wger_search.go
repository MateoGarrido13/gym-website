package web

import (
	"net/http"
	"strings"

	"PRACTICO_DOS/internal/integration"
	"PRACTICO_DOS/views"
)

// SearchRoutineExercisesHandler maneja la petición HTMX parcial desde el formulario de rutinas.
// Ruta sugerida: GET /rutinas/ejercicios/buscar
func SearchRoutineExercisesHandler(wgerACL integration.WgerACL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))

		// Validación defensiva en backend: si hay menos de 3 caracteres, no se consulta la API.
		if len([]rune(query)) < 3 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		exercises, err := wgerACL.SearchExercises(r.Context(), query, 100)
		if err != nil {
			http.Error(w, "Error consultando ejercicios externos", http.StatusBadGateway)
			return
		}

		options := make([]views.WgerExerciseOption, 0, len(exercises))
		for _, ex := range exercises {
			options = append(options, views.WgerExerciseOption{
				ID:          ex.ID,
				Name:        ex.Name,
				MuscleGroup: ex.MuscleGroup,
			})
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = views.WgerExerciseResults(options).Render(r.Context(), w)
	}
}