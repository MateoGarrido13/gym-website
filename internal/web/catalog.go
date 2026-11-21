package web

import (
	"net/http"
)

func (h *Handler) Ejercicios(w http.ResponseWriter, r *http.Request) {
	h.showCatalog(w, r, "ejercicios", catalogView{})
}

func (h *Handler) EditEjercicio(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "ejercicios")
}

func (h *Handler) CreateEjercicio(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "ejercicios", false)
}

func (h *Handler) UpdateEjercicio(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "ejercicios", true)
}

func (h *Handler) DeleteEjercicio(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "ejercicios")
}

func (h *Handler) Alumnos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	h.showCatalog(w, r, "alumnos", catalogView{})
}

func (h *Handler) CreateAlumno(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "alumnos", false)
}

func (h *Handler) EditAlumno(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "alumnos")
}

func (h *Handler) UpdateAlumno(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "alumnos", true)
}

func (h *Handler) DeleteAlumno(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "alumnos")
}

func (h *Handler) Profesores(w http.ResponseWriter, r *http.Request) {
	h.showCatalog(w, r, "profesores", catalogView{})
}

func (h *Handler) CreateProfesor(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "profesores", false)
}

func (h *Handler) EditProfesor(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "profesores")
}

func (h *Handler) UpdateProfesor(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "profesores", true)
}

func (h *Handler) DeleteProfesor(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "profesores")
}

func (h *Handler) Rutinas(w http.ResponseWriter, r *http.Request) {
	h.showCatalog(w, r, "rutinas", catalogView{})
}

func (h *Handler) EditRutina(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "rutinas")
}

func (h *Handler) CreateRutina(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "rutinas", false)
}

func (h *Handler) UpdateRutina(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "rutinas", true)
}

func (h *Handler) DeleteRutina(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "rutinas")
}

func (h *Handler) Clases(w http.ResponseWriter, r *http.Request) {
	h.showCatalog(w, r, "clases", catalogView{})
}

func (h *Handler) EditClase(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "clases")
}

func (h *Handler) CreateClase(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "clases", false)
}

func (h *Handler) UpdateClase(w http.ResponseWriter, r *http.Request) {
	h.saveCatalog(w, r, "clases", true)
}

func (h *Handler) DeleteClase(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "clases")
}
