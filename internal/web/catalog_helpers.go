package web

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"PRACTICO_DOS/internal/service"
	"github.com/lib/pq"
)

func (h *Handler) catalogRequestError(w http.ResponseWriter, r *http.Request, kind string, err error) {
	if errors.Is(err, service.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	log.Printf("no se pudo cargar el registro %s: %v", kind, err)
	http.Error(w, "No se pudo cargar el registro.", http.StatusInternalServerError)
}

func catalogID(r *http.Request) (int32, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("identificador inválido")
	}
	return int32(id), nil
}

func positiveInt32(value, message string) (int32, error) {
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil || parsed <= 0 {
		return 0, errors.New(message)
	}
	return int32(parsed), nil
}

func optionalPositiveInt32(value, message string) (*int32, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := positiveInt32(value, message)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func catalogPathForRedirect(kind string) string {
	if kind == "alumnos" {
		return "/"
	}
	return "/" + kind
}

func catalogFormError(err error) (string, int) {
	if errors.Is(err, service.ErrInvalid) {
		return err.Error(), http.StatusBadRequest
	}
	if errors.Is(err, service.ErrNotFound) {
		return "El registro ya no existe. Actualiza la página e inténtalo nuevamente.", http.StatusNotFound
	}
	if errors.Is(err, service.ErrConflict) {
		return "No se puede eliminar porque el registro está relacionado con otros datos.", http.StatusConflict
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23503" {
		return "No se pudo completar el cambio: verifica que el profesor exista y que el registro no esté asociado a otros datos.", http.StatusConflict
	}
	log.Printf("error al guardar cambios del catálogo: %v", err)
	return "No se pudo guardar el registro. Intenta nuevamente.", http.StatusInternalServerError
}

func catalogDescription(kind string) string {
	switch kind {
	case "alumnos":
		return "Registra alumnos y consulta sus datos principales."
	case "profesores":
		return "Registra compañeros de trabajo y consulta sus especialidades."
	case "ejercicios":
		return "Crea y administra el catálogo de ejercicios disponibles para las rutinas."
	case "rutinas":
		return "Define planes de entrenamiento y asígnalos a un profesor."
	case "clases":
		return "Registra clases, su capacidad y el profesor responsable."
	default:
		return ""
	}
}
