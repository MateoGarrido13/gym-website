package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"
)

func createAlumnoDePrueba(t *testing.T, queries *db.Queries, ctx context.Context) db.GetAlumnoCompletoRow {
	t.Helper()

	params := db.CreateAlumnoParams{
		Email:            fmt.Sprintf("alumno-%d@example.com", time.Now().UnixNano()),
		Contrasena:       "secret123",
		Nombre:           "Mateo",
		Apellido:         "Alumno",
		Telefono:         sql.NullString{String: "123456789", Valid: true},
		FechaInscripcion: time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC),
		FechaVto:         sql.NullTime{Time: time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		TipoPlan:         sql.NullString{String: "mensual", Valid: true},
		RutinaID:         sql.NullInt32{Valid: false},
	}

	alumno, err := queries.CreateAlumno(ctx, params)
	if err != nil {
		t.Fatalf("error al crear alumno de prueba: %v", err)
	}

	t.Cleanup(func() {
		_ = queries.DeleteUsuario(ctx, alumno.ID)
	})

	return alumno
}

func mismaFecha(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func TestCreateAlumno(t *testing.T) {
	queries, ctx := setup(t)

	fechaInscripcion := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	fechaVto := time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)
	params := db.CreateAlumnoParams{
		Email:            fmt.Sprintf("alumno-create-%d@example.com", time.Now().UnixNano()),
		Contrasena:       "secret123",
		Nombre:           "Mateo",
		Apellido:         "Alumno",
		Telefono:         sql.NullString{String: "123456789", Valid: true},
		FechaInscripcion: fechaInscripcion,
		FechaVto:         sql.NullTime{Time: fechaVto, Valid: true},
		TipoPlan:         sql.NullString{String: "mensual", Valid: true},
		RutinaID:         sql.NullInt32{Valid: false},
	}

	alumno, err := queries.CreateAlumno(ctx, params)
	if err != nil {
		t.Fatalf("CREATE de alumno falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteUsuario(ctx, alumno.ID)
	})

	if alumno.ID == 0 {
		t.Fatal("CREATE no devolvió un ID válido")
	}

	padre, err := queries.GetUsuario(ctx, alumno.ID)
	if err != nil {
		t.Fatalf("CreateAlumno no insertó el usuario padre: %v", err)
	}
	if padre.ID != alumno.ID {
		t.Errorf("id del padre = %d, se esperaba %d", padre.ID, alumno.ID)
	}
	if padre.Email != params.Email {
		t.Errorf("email del padre = %q, se esperaba %q", padre.Email, params.Email)
	}
	if padre.Nombre != params.Nombre || padre.Apellido != params.Apellido {
		t.Errorf("nombre del padre = %q %q, se esperaba %q %q", padre.Nombre, padre.Apellido, params.Nombre, params.Apellido)
	}
	if padre.Rol != "alumno" {
		t.Errorf("rol del padre = %q, se esperaba %q", padre.Rol, "alumno")
	}
	if alumno.Email != padre.Email {
		t.Errorf("email del join = %q, se esperaba %q", alumno.Email, padre.Email)
	}
	if alumno.TipoPlan.String != params.TipoPlan.String {
		t.Errorf("tipo_plan = %q, se esperaba %q", alumno.TipoPlan.String, params.TipoPlan.String)
	}
	if !mismaFecha(alumno.FechaInscripcion, fechaInscripcion) {
		t.Errorf("fecha_inscripcion = %s, se esperaba %s", alumno.FechaInscripcion.Format("2006-01-02"), fechaInscripcion.Format("2006-01-02"))
	}
}

func TestGetAlumno(t *testing.T) {
	queries, ctx := setup(t)
	creado := createAlumnoDePrueba(t, queries, ctx)

	alumno, err := queries.GetAlumnoCompleto(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ falló: %v", err)
	}

	if alumno.ID != creado.ID {
		t.Errorf("id = %d, se esperaba %d", alumno.ID, creado.ID)
	}
	if alumno.Email != creado.Email {
		t.Errorf("email = %q, se esperaba %q", alumno.Email, creado.Email)
	}
	if alumno.Nombre != creado.Nombre || alumno.Apellido != creado.Apellido {
		t.Errorf("nombre = %q %q, se esperaba %q %q", alumno.Nombre, alumno.Apellido, creado.Nombre, creado.Apellido)
	}
	if alumno.TipoPlan.String != creado.TipoPlan.String {
		t.Errorf("tipo_plan = %q, se esperaba %q", alumno.TipoPlan.String, creado.TipoPlan.String)
	}
	if !mismaFecha(alumno.FechaInscripcion, creado.FechaInscripcion) {
		t.Errorf("fecha_inscripcion = %s, se esperaba %s", alumno.FechaInscripcion.Format("2006-01-02"), creado.FechaInscripcion.Format("2006-01-02"))
	}

	padre, err := queries.GetUsuario(ctx, creado.ID)
	if err != nil {
		t.Fatalf("el usuario padre debería existir: %v", err)
	}
	if padre.Email != alumno.Email {
		t.Errorf("email del padre = %q, se esperaba %q", padre.Email, alumno.Email)
	}
}

func TestUpdateAlumno(t *testing.T) {
	queries, ctx := setup(t)
	creado := createAlumnoDePrueba(t, queries, ctx)

	nuevaFechaVto := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	updateParams := db.UpdateAlumnoDetallesParams{
		UsuarioID: creado.ID,
		FechaVto:  sql.NullTime{Time: nuevaFechaVto, Valid: true},
		TipoPlan:  sql.NullString{String: "anual", Valid: true},
		RutinaID:  sql.NullInt32{Valid: false},
	}

	if err := queries.UpdateAlumnoDetalles(ctx, updateParams); err != nil {
		t.Fatalf("UPDATE falló: %v", err)
	}

	alumno, err := queries.GetAlumnoCompleto(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ después del UPDATE falló: %v", err)
	}

	if alumno.TipoPlan.String != updateParams.TipoPlan.String {
		t.Errorf("tipo_plan = %q, se esperaba %q", alumno.TipoPlan.String, updateParams.TipoPlan.String)
	}
	if !alumno.FechaVto.Valid || !mismaFecha(alumno.FechaVto.Time, nuevaFechaVto) {
		t.Errorf("fecha_vto = %v, se esperaba %s", alumno.FechaVto, nuevaFechaVto.Format("2006-01-02"))
	}
	if alumno.Email != creado.Email {
		t.Errorf("el email del usuario padre no debió cambiar: %q", alumno.Email)
	}
}

func TestDeleteAlumno(t *testing.T) {
	queries, ctx := setup(t)
	creado := createAlumnoDePrueba(t, queries, ctx)

	if err := queries.DeleteUsuario(ctx, creado.ID); err != nil {
		t.Fatalf("DELETE del usuario padre falló: %v", err)
	}

	_, err := queries.GetAlumnoCompleto(ctx, creado.ID)
	if err == nil {
		t.Fatal("el alumno no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE del alumno: %v", err)
	}

	_, err = queries.GetUsuario(ctx, creado.ID)
	if err == nil {
		t.Fatal("el usuario padre no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE del usuario padre: %v", err)
	}
}

func TestListAlumnos(t *testing.T) {
	queries, ctx := setup(t)
	creado := createAlumnoDePrueba(t, queries, ctx)

	alumnos, err := queries.ListAlumnosCompletos(ctx)
	if err != nil {
		t.Fatalf("LIST falló: %v", err)
	}

	encontrado := false
	for _, a := range alumnos {
		if a.ID == creado.ID {
			encontrado = true
			if a.Email != creado.Email {
				t.Errorf("email = %q, se esperaba %q", a.Email, creado.Email)
			}
			if a.TipoPlan.String != creado.TipoPlan.String {
				t.Errorf("tipo_plan = %q, se esperaba %q", a.TipoPlan.String, creado.TipoPlan.String)
			}
			break
		}
	}
	if !encontrado {
		t.Fatalf("el alumno creado (ID %d) no apareció en el listado", creado.ID)
	}
}
