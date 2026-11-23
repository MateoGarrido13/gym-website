package web

import (
	"fmt"
	"log"
	"net/http"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

// catalogView agrupa lo que define cómo se dibuja el formulario y la respuesta.
// El valor cero significa "formulario de crear en blanco, sin error, status 200".
type catalogView struct {
	Values  views.FormValues // texto precargado en los inputs
	Editing bool             // true si es el formulario de edición
	ID      int32            // id del registro; solo se usa si Editing
	Message string           // mensaje de error opcional
	Status  int              // código HTTP; 0 equivale a 200
}

func (h *Handler) showCatalog(w http.ResponseWriter, r *http.Request, kind string, view catalogView) {
	status := view.Status
	if status == 0 {
		status = http.StatusOK
	}
	page, err := h.catalogPage(r, kind, view)
	if err != nil {
		log.Printf("no se pudo cargar %s para la página: %v", kind, err)
		http.Error(w, "No se pudieron cargar los registros.", http.StatusInternalServerError)
		return
	}
	h.renderStatus(w, r, page, status)
}

func (h *Handler) catalogPage(r *http.Request, kind string, view catalogView) (views.PageData, error) {
	resource, ok := h.resourceFor(kind)
	if !ok {
		return views.PageData{}, service.ErrNotFound
	}

	page := views.PageData{
		Active:      kind,
		Kind:        kind,
		Action:      "/" + kind,
		Title:       resource.Title(),
		Headers:     resource.Headers(),
		Description: resource.Description(),
		FormError:   view.Message,
		IsEditing:   view.Editing,
		ShowActions: true,
		FormValues:  view.Values,
	}

	entities, err := resource.Entities(r)
	if err != nil {
		return page, err
	}
	page.Entities = entities

	if paged, ok := resource.(pagedCatalog); ok {
		if err := paged.ApplyPage(r, &page); err != nil {
			return page, err
		}
	}

	if err := resource.FormOptions(r, &page, view.Values, view.Editing); err != nil {
		return page, err
	}

	if view.Editing {
		page.Action = fmt.Sprintf("/%s/%d", kind, view.ID)
	}
	return page, nil
}
