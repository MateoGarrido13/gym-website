package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
)

type ClaseHorarioHandler struct {
	svc *service.ClaseHorarioService
}

func NewClaseHorarioHandler(svc *service.ClaseHorarioService) *ClaseHorarioHandler {
	return &ClaseHorarioHandler{svc: svc}
}

func (h *ClaseHorarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	claseID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var in service.ClaseHorarioCreateInput
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	out, err := h.svc.Create(r.Context(), claseID, in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *ClaseHorarioHandler) List(w http.ResponseWriter, r *http.Request) {
	claseID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.svc.List(r.Context(), claseID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ClaseHorarioHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	out, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ClaseHorarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
