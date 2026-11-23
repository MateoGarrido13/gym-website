package web

import (
	"net/http"

	"PRACTICO_DOS/views"
)

// catalogResource aísla, por cada entidad del catálogo, la lógica que antes
// vivía mezclada en un switch gigante (catalogPage, createCatalog, etc). El
// único lugar que todavía elige "qué entidad es" es resourceFor.
type catalogResource interface {
	Title() string
	Headers() []string
	Description() string

	// Entities construye las filas de la tabla a partir de la base de datos.
	Entities(r *http.Request) ([]views.Entity, error)

	// FormOptions completa los datos auxiliares del formulario (combos,
	// valores por defecto). No todas las entidades necesitan algo acá.
	FormOptions(r *http.Request, page *views.PageData, values views.FormValues, editing bool) error

	// EditValues carga los valores de un registro existente para precargar
	// el formulario de edición.
	EditValues(r *http.Request, id int32) (views.FormValues, error)

	// Create da de alta un registro con los valores del formulario. Los
	// errores de validación propia deben envolverse con &formError{...} para
	// distinguirse de errores del servicio.
	Create(r *http.Request, values views.FormValues) error

	// Update modifica el registro id con los valores del formulario. Mismo
	// criterio de errores que Create.
	Update(r *http.Request, id int32, values views.FormValues) error

	Delete(r *http.Request, id int32) error
}

// pagedCatalog completa el total y la ventana limit/offset del listado.
// Solo lo implementa alumnos.
type pagedCatalog interface {
	ApplyPage(r *http.Request, page *views.PageData) error
}

// formError representa un error de validación propio del formulario, con su
// mensaje final para el usuario y el status HTTP a devolver.
type formError struct {
	message string
	status  int
}

func (e *formError) Error() string { return e.message }

func (h *Handler) resourceFor(kind string) (catalogResource, bool) {
	switch kind {
	case "alumnos":
		return alumnosResource{h}, true
	case "profesores":
		return profesoresResource{h}, true
	case "ejercicios":
		return ejerciciosResource{h}, true
	case "rutinas":
		return rutinasResource{h}, true
	case "clases":
		return clasesResource{h}, true
	default:
		return nil, false
	}
}
