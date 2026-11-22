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

// renderList vuelve a consultar "kind" en la base y renderiza SOLO EntityList,
// sin Layout ni EntityForm alrededor. Es lo que usan los handlers de alta: la
// respuesta es justo el fragmento que HTMX va a pegar en #entity-list.
func (h *Handler) renderList(w http.ResponseWriter, r *http.Request, kind string, status int) {
	page, err := h.catalogPage(r, kind, catalogView{})
	if err != nil {
		log.Printf("no se pudo volver a listar %s: %v", kind, err)
		http.Error(w, "No se pudieron cargar los registros.", http.StatusInternalServerError)
		return
	}
	var output bytes.Buffer
	if err := views.EntityList(page).Render(r.Context(), &output); err != nil {
		log.Printf("no se pudo renderizar la lista de %s: %v", kind, err)
		http.Error(w, "No se pudo renderizar la lista.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := output.WriteTo(w); err != nil {
		log.Printf("no se pudo escribir la lista de %s: %v", kind, err)
	}
}

func (h *Handler) renderStatus(w http.ResponseWriter, r *http.Request, page views.PageData, status int) {
	if user, ok := UserFrom(r.Context()); ok {
		page.UserName = user.Nombre + " " + user.Apellido
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
