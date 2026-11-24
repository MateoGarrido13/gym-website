package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "PRACTICO_DOS/db/sqlc"
	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

type alumnosResource struct{ h *Handler }

func (res alumnosResource) Title() string { return "Alumnos" }

func (res alumnosResource) Headers() []string {
	return []string{"NOMBRE", "CORREO ELECTRÓNICO", "TELÉFONO", "PLAN / VENCIMIENTO"}
}

func (res alumnosResource) Description() string {
	return "Registra alumnos y consulta sus datos principales."
}

const alumnosPorPagina int32 = 10

func alumnosWindow(r *http.Request) (int32, int32) {
	offset := int32(0)
	if raw := r.URL.Query().Get("offset"); raw != "" {
		n, err := strconv.ParseInt(raw, 10, 32)
		if err == nil && n > 0 {
			offset = int32(n)
		}
	}
	return alumnosPorPagina, offset
}

func (res alumnosResource) ApplyPage(r *http.Request, page *views.PageData) error {
	limit, offset := alumnosWindow(r)
	total, err := res.h.alumnos.Count(r.Context())
	if err != nil {
		return err
	}
	page.Paged = true
	page.Total = int(total)
	page.Limit = int(limit)
	page.Offset = int(offset)
	return nil
}

func (res alumnosResource) Entities(r *http.Request) ([]views.Entity, error) {
	limit, offset := alumnosWindow(r)
	items, err := res.h.alumnos.ListPage(r.Context(), limit, offset)
	if err != nil {
		return nil, err
	}
	entities := make([]views.Entity, 0, len(items))
	for _, item := range items {
		detail := ""
	if item.TipoPlan != nil {
		detail = string(*item.TipoPlan)
	}
		if item.FechaVto != nil {
			if detail != "" {
				detail += " · "
			}
			detail += "Vence " + *item.FechaVto
		}
		entities = append(entities, views.Entity{
			Values:    []string{item.Nombre + " " + item.Apellido, item.Email, stringValue(item.Telefono), detail},
			EditURL:   fmt.Sprintf("/alumnos/%d/editar", item.ID),
			DeleteURL: fmt.Sprintf("/alumnos/%d", item.ID),
		})
	}
	return entities, nil
}

func (res alumnosResource) FormOptions(r *http.Request, page *views.PageData, values views.FormValues, editing bool) error {
	page.FormValues.FechaInscripcion = defaultEnrollmentDate(values.FechaInscripcion)
	if !editing {
		return nil
	}
	routines, err := res.h.routineOptions(r)
	if err != nil {
		return err
	}
	page.Options = routines
	return nil
}

func (res alumnosResource) EditValues(r *http.Request, id int32) (views.FormValues, error) {
	item, err := res.h.alumnos.Get(r.Context(), id)
	if err != nil {
		return views.FormValues{}, err
	}
	var tipoPlanStr string
	if item.TipoPlan != nil {
		tipoPlanStr = string(*item.TipoPlan)
	}
	values := views.FormValues{
		TipoPlan: tipoPlanStr,
		FechaVto: stringValue(item.FechaVto),
	}
	if item.RutinaID != nil {
		values.RutinaID = strconv.FormatInt(int64(*item.RutinaID), 10)
	}
	return values, nil
}

// tipoPlanFromForm convierte el string del formulario a *db.EnumTipoPlan validando que sea un valor permitido del enum.
func tipoPlanFromForm(s string) (*db.EnumTipoPlan, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	tp, err := service.ParseTipoPlan(s)
	if err != nil {
		return nil, err
	}
	return &tp, nil
}

func (res alumnosResource) Create(r *http.Request, values views.FormValues) error {
	_, err := res.insert(r, values)
	return err
}

func (res alumnosResource) insert(r *http.Request, values views.FormValues) (service.Alumno, error) {
	if err := validateIdentity(values); err != nil {
		return service.Alumno{}, &formError{err.Error(), http.StatusBadRequest}
	}
	tipoPlan, err := tipoPlanFromForm(values.TipoPlan)
	if err != nil {
		return service.Alumno{}, &formError{err.Error(), http.StatusBadRequest}
	}

	return res.h.alumnos.Create(r.Context(), service.AlumnoCreateInput{
		Email:            values.Email,
		Contrasena:       values.Contrasena,
		Nombre:           values.Nombre,
		Apellido:         values.Apellido,
		Telefono:         optional(values.Telefono),
		FechaInscripcion: values.FechaInscripcion,
		FechaVto:         optional(values.FechaVto),
		TipoPlan:         tipoPlan,
	})
}

func (res alumnosResource) Update(r *http.Request, id int32, values views.FormValues) error {
	routineID, err := optionalPositiveInt32(values.RutinaID, "Selecciona una rutina válida.")
	if err != nil {
		return &formError{err.Error(), http.StatusBadRequest}
	}
	tipoPlan, err := tipoPlanFromForm(values.TipoPlan)
	if err != nil {
		return &formError{err.Error(), http.StatusBadRequest}
	}

	_, err = res.h.alumnos.Update(r.Context(), id, service.AlumnoUpdateInput{
		FechaVto: optional(values.FechaVto),
		TipoPlan: tipoPlan,
		RutinaID: routineID,
	})
	return err
}

func (res alumnosResource) Delete(r *http.Request, id int32) error {
	return res.h.alumnos.Delete(r.Context(), id)
}
