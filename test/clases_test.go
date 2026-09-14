package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCreateClase(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	params := db.CreateClaseParams{
		Nombre:     fmt.Sprintf("yoga-%d", time.Now().UnixNano()),
		MaxAlumnos: 20,
		ProfesorID: profesor.ID,
	}

	clase, err := queries.CreateClase(ctx, params)
	if err != nil {
		t.Fatalf("CREATE de clase falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteClase(ctx, clase.ID)
	})

	if clase.ID == 0 {
		t.Fatal("CREATE no devolvió un ID válido")
	}
	if clase.Nombre != params.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", clase.Nombre, params.Nombre)
	}
	if clase.MaxAlumnos != params.MaxAlumnos {
		t.Errorf("max_alumnos = %d, se esperaba %d", clase.MaxAlumnos, params.MaxAlumnos)
	}
	if clase.ProfesorID != profesor.ID {
		t.Errorf("profesor_id = %d, se esperaba %d", clase.ProfesorID, profesor.ID)
	}
}

func TestGetClase(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createClaseDePrueba(t, profesor.ID)

	clase, err := queries.GetClase(ctx, creada.ID)
	if err != nil {
		t.Fatalf("READ falló: %v", err)
	}
	if clase.ID != creada.ID || clase.Nombre != creada.Nombre {
		t.Errorf("clase = %+v, se esperaba %+v", clase, creada)
	}
}

func TestUpdateClase(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createClaseDePrueba(t, profesor.ID)

	updateParams := db.UpdateClaseParams{
		ID:         creada.ID,
		Nombre:     fmt.Sprintf("pilates-%d", time.Now().UnixNano()),
		MaxAlumnos: 8,
	}

	if err := queries.UpdateClase(ctx, updateParams); err != nil {
		t.Fatalf("UPDATE falló: %v", err)
	}

	clase, err := queries.GetClase(ctx, creada.ID)
	if err != nil {
		t.Fatalf("READ después del UPDATE falló: %v", err)
	}
	if clase.Nombre != updateParams.Nombre {
		t.Errorf("nombre = %q, se esperaba %q", clase.Nombre, updateParams.Nombre)
	}
	if clase.MaxAlumnos != updateParams.MaxAlumnos {
		t.Errorf("max_alumnos = %d, se esperaba %d", clase.MaxAlumnos, updateParams.MaxAlumnos)
	}
}

func TestDeleteClase(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createClaseDePrueba(t, profesor.ID)

	if err := queries.DeleteClase(ctx, creada.ID); err != nil {
		t.Fatalf("DELETE falló: %v", err)
	}

	_, err := queries.GetClase(ctx, creada.ID)
	if err == nil {
		t.Fatal("la clase no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE: %v", err)
	}
}

func TestListClases(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	creada := createClaseDePrueba(t, profesor.ID)

	clases, err := queries.ListClases(ctx)
	if err != nil {
		t.Fatalf("LIST falló: %v", err)
	}

	encontrada := false
	for _, c := range clases {
		if c.ID == creada.ID {
			encontrada = true
			break
		}
	}
	if !encontrada {
		t.Fatalf("la clase creada (ID %d) no apareció en el listado", creada.ID)
	}
}

func TestCreateClaseConProfesorInexistente(t *testing.T) {
	_, err := queries.CreateClase(ctx, db.CreateClaseParams{
		Nombre:     "sin-profesor",
		MaxAlumnos: 10,
		ProfesorID: -1,
	})
	if err == nil {
		t.Fatal("CREATE debió fallar por FK de profesor_id")
	}
	if !esCodigoPostgres(err, "23503") {
		t.Fatalf("se esperaba violación de FK (23503), se obtuvo: %v", err)
	}
}

func TestCreateYGetClaseHorario(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	hora := horaDePrueba(9, 30)

	horario, err := queries.CreateClaseHorario(ctx, db.CreateClaseHorarioParams{
		ClaseID:    clase.ID,
		DiaSemana:  "lunes",
		HoraInicio: hora,
	})
	if err != nil {
		t.Fatalf("CREATE de horario falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteClaseHorario(ctx, horario.ID)
	})

	leido, err := queries.GetClaseHorario(ctx, horario.ID)
	if err != nil {
		t.Fatalf("READ de horario falló: %v", err)
	}
	if leido.ClaseID != clase.ID {
		t.Errorf("clase_id = %d, se esperaba %d", leido.ClaseID, clase.ID)
	}
	if leido.DiaSemana != "lunes" {
		t.Errorf("dia_semana = %q, se esperaba %q", leido.DiaSemana, "lunes")
	}
	if !mismaHora(leido.HoraInicio, hora) {
		t.Errorf("hora_inicio = %s, se esperaba %s", leido.HoraInicio.Format("15:04:05"), hora.Format("15:04:05"))
	}
}

func TestListHorariosDeClase(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "martes", horaDePrueba(18, 0))

	horarios, err := queries.ListHorariosDeClase(ctx, clase.ID)
	if err != nil {
		t.Fatalf("LIST de horarios falló: %v", err)
	}

	encontrado := false
	for _, h := range horarios {
		if h.ID == horario.ID {
			encontrado = true
			break
		}
	}
	if !encontrado {
		t.Fatalf("el horario creado (ID %d) no apareció en el listado", horario.ID)
	}
}

func TestDeleteClaseHorario(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "miercoles", horaDePrueba(7, 0))

	if err := queries.DeleteClaseHorario(ctx, horario.ID); err != nil {
		t.Fatalf("DELETE de horario falló: %v", err)
	}

	_, err := queries.GetClaseHorario(ctx, horario.ID)
	if err == nil {
		t.Fatal("el horario no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE: %v", err)
	}
}

func TestCreateHorarioConDiaInvalido(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)

	_, err := queries.CreateClaseHorario(ctx, db.CreateClaseHorarioParams{
		ClaseID:    clase.ID,
		DiaSemana:  "miércoles",
		HoraInicio: horaDePrueba(10, 0),
	})
	if err == nil {
		t.Fatal("CREATE debió fallar por dia_semana inválido")
	}
	if !esCodigoPostgres(err, "22P02") {
		t.Fatalf("se esperaba invalid_text_representation (22P02), se obtuvo: %v", err)
	}
}

func TestDeleteClaseCascadaHorarios(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "jueves", horaDePrueba(19, 0))

	if err := queries.DeleteClase(ctx, clase.ID); err != nil {
		t.Fatalf("DELETE de clase falló: %v", err)
	}

	_, err := queries.GetClaseHorario(ctx, horario.ID)
	if err == nil {
		t.Fatal("el horario debió eliminarse en cascada al borrar la clase")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado: %v", err)
	}
}
