package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"PRACTICO_DOS/internal/service"

	"github.com/lib/pq"
)

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorBody{Error: msg})
}

func readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("el cuerpo debe ser un único objeto JSON")
	}
	return nil
}

func pathID(r *http.Request) (int32, error) {
	return pathValueID(r, "id")
}

func pathValueID(r *http.Request, name string) (int32, error) {
	n, err := strconv.ParseInt(r.PathValue(name), 10, 32)
	if err != nil || n <= 0 {
		return 0, errors.New("id inválido")
	}
	return int32(n), nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrInvalid) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			writeError(w, http.StatusConflict, "el recurso ya existe")
			return
		case "23503":
			writeError(w, http.StatusBadRequest, "referencia inválida")
			return
		case "23502":
			writeError(w, http.StatusBadRequest, "falta un campo obligatorio")
			return
		}
	}

	writeError(w, http.StatusInternalServerError, "error interno")
}
