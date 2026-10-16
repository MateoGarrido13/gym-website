package handler

import (
	"net/http"

	"PRACTICO_DOS/internal/service"
)

type ProfesorHandler struct {
	svc *service.ProfesorService
}

func NewProfesorHandler(svc *service.ProfesorService) *ProfesorHandler {
	return &ProfesorHandler{svc: svc}
}

func (h *ProfesorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in service.ProfesorCreateInput
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	out, err := h.svc.Create(r.Context(), in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *ProfesorHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.List(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ProfesorHandler) Get(w http.ResponseWriter, r *http.Request) {
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

func (h *ProfesorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var in service.ProfesorUpdateInput
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	out, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ProfesorHandler) Delete(w http.ResponseWriter, r *http.Request) {
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