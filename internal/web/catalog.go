package web

import (
	"log"
	"net/http"
)

func (h *Handler) Ejercicios(w http.ResponseWriter, r *http.Request) {
	h.showCatalog(w, r, "ejercicios", catalogView{})
}

func (h *Handler) EditEjercicio(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "ejercicios")
}

func (h *Handler) CreateEjercicio(w http.ResponseWriter, r *http.Request) {
	h.createCatalog(w, r, "ejercicios")
}

func (h *Handler) UpdateEjercicio(w http.ResponseWriter, r *http.Request) {
	h.updateCatalog(w, r, "ejercicios")
}

func (h *Handler) DeleteEjercicio(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "ejercicios")
}

// Home es el punto de entrada ("/"): no muestra contenido, solo redirige según
// la sesión (profesor → /alumnos, alumno → /mi-rutina, sin sesión → /login).
// GET / captura toda ruta no registrada, por eso devuelve 404 si no es "/".
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	user, ok, err := h.currentUser(r)
	if err != nil {
		log.Printf("no se pudo verificar la sesión: %v", err)
		http.Error(w, "No se pudo verificar la sesión.", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, homeFor(user.Rol), http.StatusSeeOther)
}

func (h *Handler) Alumnos(w http.ResponseWriter, r *http.Request) {
	h.showCatalog(w, r, "alumnos", catalogView{})
}

func (h *Handler) CreateAlumno(w http.ResponseWriter, r *http.Request) {
	h.createCatalog(w, r, "alumnos")
}

func (h *Handler) EditAlumno(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "alumnos")
}

func (h *Handler) UpdateAlumno(w http.ResponseWriter, r *http.Request) {
	h.updateCatalog(w, r, "alumnos")
}

func (h *Handler) DeleteAlumno(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "alumnos")
}

func (h *Handler) Profesores(w http.ResponseWriter, r *http.Request) {
	h.showCatalog(w, r, "profesores", catalogView{})
}

func (h *Handler) CreateProfesor(w http.ResponseWriter, r *http.Request) {
	h.createCatalog(w, r, "profesores")
}

func (h *Handler) EditProfesor(w http.ResponseWriter, r *http.Request) {
	h.editCatalog(w, r, "profesores")
}

func (h *Handler) UpdateProfesor(w http.ResponseWriter, r *http.Request) {
	h.updateCatalog(w, r, "profesores")
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
	h.createCatalog(w, r, "rutinas")
}

func (h *Handler) UpdateRutina(w http.ResponseWriter, r *http.Request) {
	h.updateCatalog(w, r, "rutinas")
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
	h.createCatalog(w, r, "clases")
}

func (h *Handler) UpdateClase(w http.ResponseWriter, r *http.Request) {
	h.updateCatalog(w, r, "clases")
}

func (h *Handler) DeleteClase(w http.ResponseWriter, r *http.Request) {
	h.deleteCatalog(w, r, "clases")
}
