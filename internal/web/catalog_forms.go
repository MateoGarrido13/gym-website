package web

import (
	"net/http"
	"strings"

	"PRACTICO_DOS/views"
)

func (h *Handler) editCatalog(w http.ResponseWriter, r *http.Request, kind string) {
	id, err := catalogID(r)
	if err != nil {
		http.Error(w, "El identificador no es válido.", http.StatusBadRequest)
		return
	}
	resource, ok := h.resourceFor(kind)
	if !ok {
		http.NotFound(w, r)
		return
	}
	values, err := resource.EditValues(r, id)
	if err != nil {
		h.catalogRequestError(w, r, kind, err)
		return
	}
	h.showCatalog(w, r, kind, catalogView{Values: values, Editing: true, ID: id})
}

// readCatalogForm resuelve el recurso y lee los valores del formulario. Si algo
// falla ya responde al cliente y devuelve ok=false.
func (h *Handler) readCatalogForm(w http.ResponseWriter, r *http.Request, kind string) (catalogResource, views.FormValues, bool) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario.", http.StatusBadRequest)
		return nil, views.FormValues{}, false
	}
	resource, ok := h.resourceFor(kind)
	if !ok {
		http.NotFound(w, r)
		return nil, views.FormValues{}, false
	}
	values := views.FormValues{
		Nombre:            strings.TrimSpace(r.FormValue("nombre")),
		Apellido:          strings.TrimSpace(r.FormValue("apellido")),
		Email:             strings.TrimSpace(r.FormValue("email")),
		Contrasena:        r.FormValue("contrasena"),
		Telefono:          strings.TrimSpace(r.FormValue("telefono")),
		FechaInscripcion:  strings.TrimSpace(r.FormValue("fecha_inscripcion")),
		FechaVto:          strings.TrimSpace(r.FormValue("fecha_vto")),
		TipoPlan:          strings.TrimSpace(r.FormValue("tipo_plan")),
		Especialidad:      strings.TrimSpace(r.FormValue("especialidad")),
		Descripcion:       strings.TrimSpace(r.FormValue("descripcion")),
		GrupoMuscular:     strings.TrimSpace(r.FormValue("grupo_muscular")),
		DuracionSemanas:   strings.TrimSpace(r.FormValue("duracion_semanas")),
		ProfesorID:        strings.TrimSpace(r.FormValue("profesor_id")),
		MaxAlumnos:        strings.TrimSpace(r.FormValue("max_alumnos")),
		RutinaID:          strings.TrimSpace(r.FormValue("rutina_id")),
		SelectedExercises: append([]string(nil), r.Form["ejercicio_id"]...),
	}
	return resource, values, true
}

// showCatalogFormError vuelve a mostrar el formulario con el mensaje que
// corresponda al error (de validación propia o del servicio).
func (h *Handler) showCatalogFormError(w http.ResponseWriter, r *http.Request, kind string, view catalogView, err error) {
	if ferr, isFormError := err.(*formError); isFormError {
		view.Message, view.Status = ferr.message, ferr.status
	} else {
		view.Message, view.Status = catalogFormError(err)
	}
	h.showCatalog(w, r, kind, view)
}

func (h *Handler) createCatalog(w http.ResponseWriter, r *http.Request, kind string) {
	resource, values, ok := h.readCatalogForm(w, r, kind)
	if !ok {
		return
	}
	if err := resource.Create(r, values); err != nil {
		h.showCatalogFormError(w, r, kind, catalogView{Values: values}, err)
		return
	}
	// HTMX: en el alta, devolvemos el fragmento de la lista actualizado en vez de redirigir.
	h.renderList(w, r, kind, http.StatusOK)
}

func (h *Handler) updateCatalog(w http.ResponseWriter, r *http.Request, kind string) {
	id, err := catalogID(r)
	if err != nil {
		http.Error(w, "El identificador no es válido.", http.StatusBadRequest)
		return
	}
	resource, values, ok := h.readCatalogForm(w, r, kind)
	if !ok {
		return
	}
	if err := resource.Update(r, id, values); err != nil {
		h.showCatalogFormError(w, r, kind, catalogView{Values: values, Editing: true, ID: id}, err)
		return
	}
	// La edición sigue siendo otra página; no la tocamos para HTMX.
	http.Redirect(w, r, catalogPathForRedirect(kind), http.StatusSeeOther)
}

func (h *Handler) deleteCatalog(w http.ResponseWriter, r *http.Request, kind string) {
	id, err := catalogID(r)
	if err != nil {
		http.Error(w, "El identificador no es válido.", http.StatusBadRequest)
		return
	}
	resource, ok := h.resourceFor(kind)
	if !ok {
		http.NotFound(w, r)
		return
	}

	if err := resource.Delete(r, id); err != nil {
		message, status := catalogFormError(err)
		h.showCatalog(w, r, kind, catalogView{Message: message, Status: status})
		return
	}
	// HTMX hizo hx-delete con hx-target="closest tr" y hx-swap="outerHTML":
	// un 200 con cuerpo vacío reemplaza esa fila por "nada", es decir la borra de la vista.
	w.WriteHeader(http.StatusOK)
}
