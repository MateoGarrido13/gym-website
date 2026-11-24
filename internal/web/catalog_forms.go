package web

import (
	"fmt"
	"log"
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
// corresponda al error (de validación propia o del servicio). La usa la
// edición, que sigue siendo una página completa, no HTMX.
func (h *Handler) showCatalogFormError(w http.ResponseWriter, r *http.Request, kind string, view catalogView, err error) {
	view.Message, view.Status = formMessage(err)
	h.showCatalog(w, r, kind, view)
}

func (h *Handler) createCatalog(w http.ResponseWriter, r *http.Request, kind string) {
	resource, values, ok := h.readCatalogForm(w, r, kind)
	if !ok {
		return
	}
	if kind == "alumnos" {
		h.createAlumno(w, r, values)
		return
	}
	if err := resource.Create(r, values); err != nil {
		// HTMX: solo el mensaje de error (hx-swap-oob), nunca la página completa.
		h.renderCreateError(w, r, err)
		return
	}
	// Tras crear exitosamente, renderizar el formulario vacío (valores por defecto)
	// con un mensaje de éxito, listo para cargar un nuevo registro.
	// El form sigue apuntando a #entity-rows. En un request HTMX pedimos
	// reemplazar el contenido del body y quedarnos en la misma URL.
	// No enviamos HX-Reselect: body: HTMX ya descarta la etiqueta body de una
	// página completa, y seleccionarla otra vez no matchea nada y deja la página en blanco.
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Retarget", "body")
		w.Header().Set("HX-Reswap", "innerHTML")
	}
	h.showCatalog(w, r, kind, catalogView{
		Values:         views.FormValues{},
		SuccessMessage: "Registro creado correctamente",
		Status:         http.StatusOK,
	})
}

// createAlumno redirige a la página donde el nuevo queda según
// (nombre, apellido, id). Appendear la fila al final rompería ese orden.
func (h *Handler) createAlumno(w http.ResponseWriter, r *http.Request, values views.FormValues) {
	alumno, err := (alumnosResource{h}).insert(r, values)
	if err != nil {
		h.renderCreateError(w, r, err)
		return
	}
	antes, err := h.alumnos.CountAntes(r.Context(), alumno.ID)
	if err != nil {
		log.Printf("no se pudo ubicar la página del alumno %d: %v", alumno.ID, err)
		http.Error(w, "No se pudo cargar el registro nuevo.", http.StatusInternalServerError)
		return
	}
	offset := (antes / alumnosPorPagina) * alumnosPorPagina
	loc := fmt.Sprintf("/alumnos?limit=%d&offset=%d", alumnosPorPagina, offset)
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", loc)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, loc, http.StatusSeeOther)
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
