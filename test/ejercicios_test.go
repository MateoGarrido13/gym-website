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

func createEjercicioDePrueba(t *testing.T, queries *db.Queries, ctx context.Context) db.Ejercicio {
	t.Helper()

	// creamos un struct con los parametros para el createEjercicio
	params := db.CreateEjercicioParams{
		Nombre:      fmt.Sprintf("sentadilla-%d", time.Now().UnixNano()),
		Descripcion: sql.NullString{String: "Ejercicio de piernas", Valid: true},
	}

	ejercicio, err := queries.CreateEjercicio(ctx, params) // pasamos el contexto e intentamos una insercion con los parametros
	if err != nil {
		t.Fatalf("error al crear ejercicio de prueba: %v", err)
	}

	t.Cleanup(func() { // limpiamos la base de datos despues de la prueba
		_ = queries.DeleteEjercicio(ctx, ejercicio.ID)
	})

	return ejercicio
}

func TestCreateEjercicio(t *testing.T) {
	queries, ctx := setup(t)

	params := db.CreateEjercicioParams{ // creamos un struct con los parametros para el createEjercicio
		Nombre:      fmt.Sprintf("press-banca-%d", time.Now().UnixNano()),
		Descripcion: sql.NullString{String: "Ejercicio de pecho", Valid: true},
	}

	ejercicio, err := queries.CreateEjercicio(ctx, params)
	if err != nil {
		t.Fatalf("CREATE falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteEjercicio(ctx, ejercicio.ID)
	})

	if ejercicio.ID == 0 {
		t.Fatal("CREATE no devolvió un ID válido")
	}
	if ejercicio.Nombre != params.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", ejercicio.Nombre, params.Nombre)
	}
	if ejercicio.Descripcion.String != params.Descripcion.String {
		t.Errorf("descripcion = %q, se esperaba %q", ejercicio.Descripcion.String, params.Descripcion.String)
	}
}

func TestGetEjercicio(t *testing.T) {
	queries, ctx := setup(t)
	creado := createEjercicioDePrueba(t, queries, ctx)

	ejercicio, err := queries.GetEjercicio(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ falló: %v", err)
	}

	if ejercicio.ID != creado.ID {
		t.Errorf("id = %d, se esperaba %d", ejercicio.ID, creado.ID)
	}
	if ejercicio.Nombre != creado.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", ejercicio.Nombre, creado.Nombre)
	}
	if ejercicio.Descripcion.String != creado.Descripcion.String {
		t.Errorf("descripcion = %q, se esperaba %q", ejercicio.Descripcion.String, creado.Descripcion.String)
	}
}

func TestUpdateEjercicio(t *testing.T) {
	queries, ctx := setup(t)
	creado := createEjercicioDePrueba(t, queries, ctx)

	updateParams := db.UpdateEjercicioParams{
		ID:          creado.ID,
		Nombre:      fmt.Sprintf("peso-muerto-%d", time.Now().UnixNano()),
		Descripcion: sql.NullString{String: "Ejercicio de espalda y posterior", Valid: true},
	}

	if err := queries.UpdateEjercicio(ctx, updateParams); err != nil {
		t.Fatalf("UPDATE falló: %v", err)
	}

	ejercicio, err := queries.GetEjercicio(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ después del UPDATE falló: %v", err)
	}

	if ejercicio.Nombre != updateParams.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", ejercicio.Nombre, updateParams.Nombre)
	}
	if ejercicio.Descripcion.String != updateParams.Descripcion.String {
		t.Errorf("descripcion = %q, se esperaba %q", ejercicio.Descripcion.String, updateParams.Descripcion.String)
	}
}

func TestDeleteEjercicio(t *testing.T) {
	queries, ctx := setup(t)
	creado := createEjercicioDePrueba(t, queries, ctx)

	if err := queries.DeleteEjercicio(ctx, creado.ID); err != nil {
		t.Fatalf("DELETE falló: %v", err)
	}

	_, err := queries.GetEjercicio(ctx, creado.ID)
	if err == nil {
		t.Fatal("el ejercicio no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) { // verificamos que el error sea ErrNoRows por haber eliminado el ejercicio
		t.Fatalf("error inesperado al verificar el DELETE: %v", err)
	}
}

func TestListEjercicios(t *testing.T) {
	queries, ctx := setup(t)
	creado := createEjercicioDePrueba(t, queries, ctx)

	ejercicios, err := queries.ListEjercicios(ctx)
	if err != nil {
		t.Fatalf("LIST falló: %v", err)
	}

	encontrado := false
	for _, e := range ejercicios {
		if e.ID == creado.ID {
			encontrado = true
			if e.Nombre != creado.Nombre {
				t.Errorf("nombre = %q, se esperaba %q", e.Nombre, creado.Nombre)
			}
			break
		}
	}
	if !encontrado {
		t.Fatalf("el ejercicio creado (ID %d) no apareció en el listado", creado.ID)
	}
}
