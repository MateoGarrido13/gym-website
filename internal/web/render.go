package web

import (
	"bytes"
	"log"
	"net/http"

	"PRACTICO_DOS/views"
)

func (h *Handler) render(w http.ResponseWriter, r *http.Request, page views.PageData) {
	h.renderStatus(w, r, page, http.StatusOK)
}

// renderAlumnosList devuelve solo el panel de la lista (paginado HTMX).
func (h *Handler) renderAlumnosList(w http.ResponseWriter, r *http.Request) {
	page, err := h.catalogPage(r, "alumnos", catalogView{})
	if err != nil {
		log.Printf("no se pudo cargar la lista de alumnos: %v", err)
		http.Error(w, "No se pudieron cargar los registros.", http.StatusInternalServerError)
		return
	}
	var output bytes.Buffer
	if err := views.AlumnosList(page).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar la lista de alumnos: %v", err)
		http.Error(w, "No se pudo renderizar la lista.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := output.WriteTo(w); err != nil {
		log.Printf("no se pudo escribir la lista de alumnos: %v", err)
	}
}

// renderNewRow renderiza la fila del registro recién creado, más dos
// fragmentos hx-swap-oob en la misma respuesta: el contador actualizado y
// la orden de borrar la fila de "vacío" (si ya no está, no hace nada).
// Vuelve a listar el catálogo, ordenado por id, y toma la última entidad,
// que es la nueva. Alumnos no usa este camino: su alta redirige a la página
// del orden (nombre, apellido, id).
func (h *Handler) renderNewRow(w http.ResponseWriter, r *http.Request, kind string) {
	resource, ok := h.resourceFor(kind)
	if !ok {
		http.NotFound(w, r)
		return
	}
	entities, err := resource.Entities(r)
	if err != nil || len(entities) == 0 {
		log.Printf("no se pudo obtener la fila nueva de %s: %v", kind, err)
		http.Error(w, "No se pudo cargar el registro nuevo.", http.StatusInternalServerError)
		return
	}
	var output bytes.Buffer
	if err := views.EntityRow(entities[len(entities)-1], true).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar la fila nueva de %s: %v", kind, err)
		http.Error(w, "No se pudo renderizar el registro nuevo.", http.StatusInternalServerError)
		return
	}
	if err := views.ResultCount(len(entities), true).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar el contador de %s: %v", kind, err)
		http.Error(w, "No se pudo renderizar el registro nuevo.", http.StatusInternalServerError)
		return
	}
	if err := views.ClearEmptyRow().Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar la limpieza de %s: %v", kind, err)
		http.Error(w, "No se pudo renderizar el registro nuevo.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := output.WriteTo(w); err != nil {
		log.Printf("no se pudo escribir la fila nueva de %s: %v", kind, err)
	}
}

// renderCreateError devuelve SOLO el mensaje de error del formulario, como
// fragmento hx-swap-oob: HTMX lo pega en el formulario sin tocar la lista
// ni recibir la página completa (que rompería el fragmento esperado).
func (h *Handler) renderCreateError(w http.ResponseWriter, r *http.Request, err error) {
	message, status := formMessage(err)
	var output bytes.Buffer
	if rErr := views.FormError(message, true).Render(r.Context(), &output); rErr != nil {
		log.Printf("no se pudo renderizar el error del formulario: %v", rErr)
		http.Error(w, "No se pudo guardar el registro.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, wErr := output.WriteTo(w); wErr != nil {
		log.Printf("no se pudo escribir el error del formulario: %v", wErr)
	}
}

func (h *Handler) renderStatus(w http.ResponseWriter, r *http.Request, page views.PageData, status int) {
	if user, ok := UserFrom(r.Context()); ok {
		page.UserName = user.Nombre + " " + user.Apellido
		page.UserRol = user.Rol
	}
	var output bytes.Buffer
	if err := views.IndexPage(page).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar la página %q: %v", page.Title, err)
		http.Error(w, "No se pudo renderizar la página.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := output.WriteTo(w); err != nil {
		log.Printf("no se pudo escribir la página %q: %v", page.Title, err)
	}
}

func (h *Handler) renderLogin(w http.ResponseWriter, r *http.Request, data views.LoginData, status int) {
	var output bytes.Buffer
	if err := views.LoginPage(data).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar el login: %v", err)
		http.Error(w, "No se pudo renderizar la página.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := output.WriteTo(w); err != nil {
		log.Printf("no se pudo escribir el login: %v", err)
	}
}

func (h *Handler) renderAlumnoRutina(w http.ResponseWriter, r *http.Request, data views.RutinaAlumnoData) {
	var output bytes.Buffer
	if err := views.AlumnoRutinaPage(data).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar la rutina del alumno: %v", err)
		http.Error(w, "No se pudo renderizar la página.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := output.WriteTo(w); err != nil {
		log.Printf("no se pudo escribir la rutina del alumno: %v", err)
	}
}
