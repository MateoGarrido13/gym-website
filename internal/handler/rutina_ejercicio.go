package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
)

type RutinaEjercicioHandler struct {
	svc *service.RutinaEjercicioService
}

func NewRutinaEjercicioHandler(svc *service.RutinaEjercicioService) *RutinaEjercicioHandler {
	return &RutinaEjercicioHandler{svc: svc}
}

func (h *RutinaEjercicioHandler) Create(w http.ResponseWriter, r *http.Request) {
	rutinaID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var in service.RutinaEjercicioCreateInput
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	out, err := h.svc.Create(r.Context(), rutinaID, in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *RutinaEjercicioHandler) List(w http.ResponseWriter, r *http.Request) {
	rutinaID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.svc.List(r.Context(), rutinaID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *RutinaEjercicioHandler) Update(w http.ResponseWriter, r *http.Request) {
	rutinaID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ejercicioID, err := pathValueID(r, "ejercicio_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var in service.RutinaEjercicioUpdateInput
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	out, err := h.svc.Update(r.Context(), rutinaID, ejercicioID, in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *RutinaEjercicioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	rutinaID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ejercicioID, err := pathValueID(r, "ejercicio_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(r.Context(), rutinaID, ejercicioID); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
