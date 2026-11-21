package web

import (
	"fmt"
	"log"
	"net/http"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

func (h *Handler) showCatalog(w http.ResponseWriter, r *http.Request, kind string, values views.FormValues, editing bool, id int32, message string, status int) {
	page, err := h.catalogPage(r, kind, values, editing, id, message)
	if err != nil {
		log.Printf("no se pudo cargar %s para la página: %v", kind, err)
		http.Error(w, "No se pudieron cargar los registros.", http.StatusInternalServerError)
		return
	}
	h.renderStatus(w, r, page, status)
}

func (h *Handler) catalogPage(r *http.Request, kind string, values views.FormValues, editing bool, id int32, message string) (views.PageData, error) {
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
		FormError:   message,
		IsEditing:   editing,
		ShowActions: true,
		FormValues:  values,
	}

	entities, err := resource.Entities(r)
	if err != nil {
		return page, err
	}
	page.Entities = entities

	if err := resource.FormOptions(r, &page, values, editing); err != nil {
		return page, err
	}

	if editing {
		page.Action = fmt.Sprintf("/%s/%d", kind, id)
	}
	return page, nil
}
