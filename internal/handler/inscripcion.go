package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
)

type InscripcionHandler struct {
	svc *service.InscripcionService
}

func NewInscripcionHandler(svc *service.InscripcionService) *InscripcionHandler {
	return &InscripcionHandler{svc: svc}
}

func (h *InscripcionHandler) Create(w http.ResponseWriter, r *http.Request) {
	horarioID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var in service.InscripcionCreateInput
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	out, err := h.svc.Create(r.Context(), horarioID, in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *InscripcionHandler) List(w http.ResponseWriter, r *http.Request) {
	horarioID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.svc.List(r.Context(), horarioID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *InscripcionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	horarioID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	alumnoID, err := pathValueID(r, "alumno_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(r.Context(), horarioID, alumnoID); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
