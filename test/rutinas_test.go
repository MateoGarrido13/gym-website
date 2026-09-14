package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCreateRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	params := db.CreateRutinaParams{
		Nombre:          fmt.Sprintf("full-body-%d", time.Now().UnixNano()),
		DuracionSemanas: 6,
		ProfesorID:      profesor.ID,
	}

	rutina, err := queries.CreateRutina(ctx, params)
	if err != nil {
		t.Fatalf("CREATE de rutina falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteRutina(ctx, rutina.ID)
	})

	if rutina.ID == 0 {
		t.Fatal("CREATE no devolvió un ID válido")
	}
	if rutina.Nombre != params.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", rutina.Nombre, params.Nombre)
	}
	if rutina.DuracionSemanas != params.DuracionSemanas {
		t.Errorf("duracion_semanas = %d, se esperaba %d", rutina.DuracionSemanas, params.DuracionSemanas)
	}
	if rutina.ProfesorID != profesor.ID {
		t.Errorf("profesor_id = %d, se esperaba %d", rutina.ProfesorID, profesor.ID)
	}
}

func TestGetRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createRutinaDePrueba(t, profesor.ID)

	rutina, err := queries.GetRutina(ctx, creada.ID)
	if err != nil {
		t.Fatalf("READ falló: %v", err)
	}
	if rutina.ID != creada.ID {
		t.Errorf("id = %d, se esperaba %d", rutina.ID, creada.ID)
	}
	if rutina.Nombre != creada.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", rutina.Nombre, creada.Nombre)
	}
}

func TestUpdateRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createRutinaDePrueba(t, profesor.ID)

	updateParams := db.UpdateRutinaParams{
		ID:              creada.ID,
		Nombre:          fmt.Sprintf("actualizada-%d", time.Now().UnixNano()),
		DuracionSemanas: 12,
	}

	if err := queries.UpdateRutina(ctx, updateParams); err != nil {
		t.Fatalf("UPDATE falló: %v", err)
	}

	rutina, err := queries.GetRutina(ctx, creada.ID)
	if err != nil {
		t.Fatalf("READ después del UPDATE falló: %v", err)
	}
	if rutina.Nombre != updateParams.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", rutina.Nombre, updateParams.Nombre)
	}
	if rutina.DuracionSemanas != updateParams.DuracionSemanas {
		t.Errorf("duracion_semanas = %d, se esperaba %d", rutina.DuracionSemanas, updateParams.DuracionSemanas)
	}
}

func TestDeleteRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createRutinaDePrueba(t, profesor.ID)

	if err := queries.DeleteRutina(ctx, creada.ID); err != nil {
		t.Fatalf("DELETE falló: %v", err)
	}

	_, err := queries.GetRutina(ctx, creada.ID)
	if err == nil {
		t.Fatal("la rutina no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE: %v", err)
	}
}

func TestListRutinas(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createRutinaDePrueba(t, profesor.ID)

	rutinas, err := queries.ListRutinas(ctx)
	if err != nil {
		t.Fatalf("LIST falló: %v", err)
	}

	encontrada := false
	for _, r := range rutinas {
		if r.ID == creada.ID {
			encontrada = true
			break
		}
	}
	if !encontrada {
		t.Fatalf("la rutina creada (ID %d) no apareció en el listado", creada.ID)
	}
}

func TestCreateRutinaConProfesorInexistente(t *testing.T) {
	_, err := queries.CreateRutina(ctx, db.CreateRutinaParams{
		Nombre:          "sin-profesor",
		DuracionSemanas: 4,
		ProfesorID:      -1,
	})
	if err == nil {
		t.Fatal("CREATE debió fallar por FK de profesor_id")
	}
	if !esCodigoPostgres(err, "23503") {
		t.Fatalf("se esperaba violación de FK (23503), se obtuvo: %v", err)
	}
}

func TestAddYListEjerciciosRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)
	ejercicio := createEjercicioDePrueba(t)

	err := queries.AddEjercicioARutina(ctx, db.AddEjercicioARutinaParams{
		RutinaID:         rutina.ID,
		EjercicioID:      ejercicio.ID,
		Orden:            1,
		Series:           sql.NullInt32{Int32: 4, Valid: true},
		Repeticiones:     sql.NullInt32{Int32: 10, Valid: true},
		DescansoSegundos: sql.NullInt32{Int32: 90, Valid: true},
	})
	if err != nil {
		t.Fatalf("AddEjercicioARutina falló: %v", err)
	}

	ejercicios, err := queries.ListEjerciciosRutina(ctx, rutina.ID)
	if err != nil {
		t.Fatalf("ListEjerciciosRutina falló: %v", err)
	}
	if len(ejercicios) != 1 {
		t.Fatalf("se esperaba 1 ejercicio, se obtuvo %d", len(ejercicios))
	}
	if ejercicios[0].ID != ejercicio.ID {
		t.Errorf("ejercicio_id = %d, se esperaba %d", ejercicios[0].ID, ejercicio.ID)
	}
	if ejercicios[0].Nombre != ejercicio.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", ejercicios[0].Nombre, ejercicio.Nombre)
	}
}

func TestEjercicioDuplicadoEnRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)
	ejercicio := createEjercicioDePrueba(t)

	params := db.AddEjercicioARutinaParams{
		RutinaID:    rutina.ID,
		EjercicioID: ejercicio.ID,
		Orden:       1,
	}
	if err := queries.AddEjercicioARutina(ctx, params); err != nil {
		t.Fatalf("primer AddEjercicioARutina falló: %v", err)
	}

	err := queries.AddEjercicioARutina(ctx, params)
	if err == nil {
		t.Fatal("no se debió poder agregar el mismo ejercicio dos veces")
	}
	if !esCodigoPostgres(err, "23505") {
		t.Fatalf("se esperaba unique violation (23505), se obtuvo: %v", err)
	}
}

func TestDeleteEjercicioDeRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)
	ejercicio := createEjercicioDePrueba(t)

	if err := queries.AddEjercicioARutina(ctx, db.AddEjercicioARutinaParams{
		RutinaID:    rutina.ID,
		EjercicioID: ejercicio.ID,
		Orden:       1,
	}); err != nil {
		t.Fatalf("AddEjercicioARutina falló: %v", err)
	}

	if err := queries.DeleteEjercicioDeRutina(ctx, db.DeleteEjercicioDeRutinaParams{
		RutinaID:    rutina.ID,
		EjercicioID: ejercicio.ID,
	}); err != nil {
		t.Fatalf("DeleteEjercicioDeRutina falló: %v", err)
	}

	ejercicios, err := queries.ListEjerciciosRutina(ctx, rutina.ID)
	if err != nil {
		t.Fatalf("ListEjerciciosRutina falló: %v", err)
	}
	if len(ejercicios) != 0 {
		t.Fatalf("la rutina debió quedar sin ejercicios, tiene %d", len(ejercicios))
	}
}

func TestDeleteRutinaCascadaEjercicios(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)
	ejercicio := createEjercicioDePrueba(t)

	if err := queries.AddEjercicioARutina(ctx, db.AddEjercicioARutinaParams{
		RutinaID:    rutina.ID,
		EjercicioID: ejercicio.ID,
		Orden:       1,
	}); err != nil {
		t.Fatalf("AddEjercicioARutina falló: %v", err)
	}

	if err := queries.DeleteRutina(ctx, rutina.ID); err != nil {
		t.Fatalf("DELETE de rutina falló: %v", err)
	}

	ejercicios, err := queries.ListEjerciciosRutina(ctx, rutina.ID)
	if err != nil {
		t.Fatalf("ListEjerciciosRutina falló: %v", err)
	}
	if len(ejercicios) != 0 {
		t.Fatal("borrar la rutina debió eliminar sus rutina_ejercicio")
	}

	if _, err := queries.GetEjercicio(ctx, ejercicio.ID); err != nil {
		t.Fatalf("el ejercicio catálogo no debió borrarse: %v", err)
	}
}

func TestDeleteEjercicioUsadoEnRutina(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)
	ejercicio := createEjercicioDePrueba(t)

	if err := queries.AddEjercicioARutina(ctx, db.AddEjercicioARutinaParams{
		RutinaID:    rutina.ID,
		EjercicioID: ejercicio.ID,
		Orden:       1,
	}); err != nil {
		t.Fatalf("AddEjercicioARutina falló: %v", err)
	}

	err := queries.DeleteEjercicio(ctx, ejercicio.ID)
	if err == nil {
		t.Fatal("no se debió poder borrar un ejercicio referenciado por una rutina")
	}
	if !esCodigoPostgres(err, "23503") {
		t.Fatalf("se esperaba violación de FK (23503), se obtuvo: %v", err)
	}
}
