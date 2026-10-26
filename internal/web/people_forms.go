package web

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"PRACTICO_DOS/internal/service"
	"github.com/lib/pq"
)

type submittedValues struct {
	Nombre           string
	Apellido         string
	Email            string
	Contrasena       string
	Telefono         string
	FechaInscripcion string
	FechaVto         string
	TipoPlan         string
	Especialidad     string
}

func formValues(r *http.Request) submittedValues {
	return submittedValues{
		Nombre:           strings.TrimSpace(r.FormValue("nombre")),
		Apellido:         strings.TrimSpace(r.FormValue("apellido")),
		Email:            strings.TrimSpace(r.FormValue("email")),
		Contrasena:       r.FormValue("contrasena"),
		Telefono:         strings.TrimSpace(r.FormValue("telefono")),
		FechaInscripcion: strings.TrimSpace(r.FormValue("fecha_inscripcion")),
		FechaVto:         strings.TrimSpace(r.FormValue("fecha_vto")),
		TipoPlan:         strings.TrimSpace(r.FormValue("tipo_plan")),
		Especialidad:     strings.TrimSpace(r.FormValue("especialidad")),
	}
}

func validateIdentity(values submittedValues) error {
	if values.Nombre == "" || values.Apellido == "" || values.Email == "" || values.Contrasena == "" {
		return errors.New("Nombre, apellido, correo y contraseña son obligatorios.")
	}
	if _, err := mail.ParseAddress(values.Email); err != nil {
		return errors.New("Ingresa un correo electrónico válido.")
	}
	return nil
}

func formError(err error) (string, int) {
	if errors.Is(err, service.ErrInvalid) {
		return err.Error(), http.StatusBadRequest
	}
	if errors.Is(err, service.ErrConflict) {
		return "Ya existe una persona registrada con ese correo electrónico.", http.StatusConflict
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return "Ya existe una persona registrada con ese correo electrónico.", http.StatusConflict
	}
	log.Printf("error al guardar datos del formulario: %v", err)
	return "No se pudo guardar el registro. Intenta nuevamente.", http.StatusInternalServerError
}
