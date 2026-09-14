package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCreateProfesor(t *testing.T) {
	params := db.CreateProfesorParams{
		Email:        fmt.Sprintf("profesor-create-%d@example.com", time.Now().UnixNano()),
		Contrasena:   "secret123",
		Nombre:       "Laura",
		Apellido:     "Profesora",
		Telefono:     sql.NullString{String: "111111111", Valid: true},
		Especialidad: sql.NullString{String: "funcional", Valid: true},
	}

	profesor, err := queries.CreateProfesor(ctx, params)
	if err != nil {
		t.Fatalf("CREATE de profesor falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteUsuario(ctx, profesor.ID)
	})

	if profesor.ID == 0 {
		t.Fatal("CREATE no devolvió un ID válido")
	}

	padre, err := queries.GetUsuario(ctx, profesor.ID)
	if err != nil {
		t.Fatalf("CreateProfesor no insertó el usuario padre: %v", err)
	}
	if padre.Rol != "profesor" {
		t.Errorf("rol del padre = %q, se esperaba %q", padre.Rol, "profesor")
	}
	if padre.Email != params.Email {
		t.Errorf("email del padre = %q, se esperaba %q", padre.Email, params.Email)
	}
	if profesor.Especialidad.String != params.Especialidad.String {
		t.Errorf("especialidad = %q, se esperaba %q", profesor.Especialidad.String, params.Especialidad.String)
	}
}

func TestGetProfesor(t *testing.T) {
	creado := createProfesorDePrueba(t)

	profesor, err := queries.GetProfesorCompleto(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ falló: %v", err)
	}

	if profesor.ID != creado.ID {
		t.Errorf("id = %d, se esperaba %d", profesor.ID, creado.ID)
	}
	if profesor.Email != creado.Email {
		t.Errorf("email = %q, se esperaba %q", profesor.Email, creado.Email)
	}
	if profesor.Especialidad.String != creado.Especialidad.String {
		t.Errorf("especialidad = %q, se esperaba %q", profesor.Especialidad.String, creado.Especialidad.String)
	}
}

func TestUpdateProfesor(t *testing.T) {
	creado := createProfesorDePrueba(t)

	updateParams := db.UpdateProfesorParams{
		UsuarioID:    creado.ID,
		Especialidad: sql.NullString{String: "hipertrofia", Valid: true},
	}

	if err := queries.UpdateProfesor(ctx, updateParams); err != nil {
		t.Fatalf("UPDATE falló: %v", err)
	}

	profesor, err := queries.GetProfesorCompleto(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ después del UPDATE falló: %v", err)
	}
	if profesor.Especialidad.String != updateParams.Especialidad.String {
		t.Errorf("especialidad = %q, se esperaba %q", profesor.Especialidad.String, updateParams.Especialidad.String)
	}
	if profesor.Email != creado.Email {
		t.Errorf("el email del usuario padre no debió cambiar: %q", profesor.Email)
	}
}

func TestDeleteProfesor(t *testing.T) {
	creado := createProfesorDePrueba(t)

	if err := queries.DeleteUsuario(ctx, creado.ID); err != nil {
		t.Fatalf("DELETE del usuario padre falló: %v", err)
	}

	_, err := queries.GetProfesorCompleto(ctx, creado.ID)
	if err == nil {
		t.Fatal("el profesor no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE del profesor: %v", err)
	}

	_, err = queries.GetUsuario(ctx, creado.ID)
	if err == nil {
		t.Fatal("el usuario padre no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE del usuario padre: %v", err)
	}
}

func TestListProfesores(t *testing.T) {
	creado := createProfesorDePrueba(t)

	profesores, err := queries.ListProfesoresCompletos(ctx)
	if err != nil {
		t.Fatalf("LIST falló: %v", err)
	}

	encontrado := false
	for _, p := range profesores {
		if p.ID == creado.ID {
			encontrado = true
			if p.Especialidad.String != creado.Especialidad.String {
				t.Errorf("especialidad = %q, se esperaba %q", p.Especialidad.String, creado.Especialidad.String)
			}
			break
		}
	}
	if !encontrado {
		t.Fatalf("el profesor creado (ID %d) no apareció en el listado", creado.ID)
	}
}

func TestGetProfesorInexistente(t *testing.T) {
	_, err := queries.GetProfesorCompleto(ctx, -1)
	if err == nil {
		t.Fatal("GetProfesorCompleto debió fallar con un id inexistente")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado: %v", err)
	}
}
