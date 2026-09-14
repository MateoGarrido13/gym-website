package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"database/sql"
	"errors"
	"testing"
)

func TestCreateInscripcion(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "viernes", horaDePrueba(17, 0))
	alumno := createAlumnoDePrueba(t)

	inscripcion, err := queries.CreateInscripcion(ctx, db.CreateInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       alumno.ID,
	})
	if err != nil {
		t.Fatalf("CREATE de inscripción falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteInscripcion(ctx, db.DeleteInscripcionParams{
			ClaseHorarioID: horario.ID,
			AlumnoID:       alumno.ID,
		})
	})

	if inscripcion.ClaseHorarioID != horario.ID {
		t.Errorf("clase_horario_id = %d, se esperaba %d", inscripcion.ClaseHorarioID, horario.ID)
	}
	if inscripcion.AlumnoID != alumno.ID {
		t.Errorf("alumno_id = %d, se esperaba %d", inscripcion.AlumnoID, alumno.ID)
	}
	if !inscripcion.FechaInscripcion.Valid {
		t.Error("fecha_inscripcion debería estar definida")
	}
}

func TestListInscriptos(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "sabado", horaDePrueba(10, 0))
	alumno := createAlumnoDePrueba(t)

	if _, err := queries.CreateInscripcion(ctx, db.CreateInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       alumno.ID,
	}); err != nil {
		t.Fatalf("CREATE de inscripción falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteInscripcion(ctx, db.DeleteInscripcionParams{
			ClaseHorarioID: horario.ID,
			AlumnoID:       alumno.ID,
		})
	})

	inscriptos, err := queries.ListInscriptos(ctx, horario.ID)
	if err != nil {
		t.Fatalf("LIST de inscriptos falló: %v", err)
	}

	encontrado := false
	for _, i := range inscriptos {
		if i.AlumnoID == alumno.ID {
			encontrado = true
			if i.Nombre != alumno.Nombre || i.Apellido != alumno.Apellido {
				t.Errorf("nombre = %q %q, se esperaba %q %q", i.Nombre, i.Apellido, alumno.Nombre, alumno.Apellido)
			}
			break
		}
	}
	if !encontrado {
		t.Fatalf("el alumno %d no apareció entre los inscriptos", alumno.ID)
	}
}

func TestDeleteInscripcion(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "domingo", horaDePrueba(11, 0))
	alumno := createAlumnoDePrueba(t)

	if _, err := queries.CreateInscripcion(ctx, db.CreateInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       alumno.ID,
	}); err != nil {
		t.Fatalf("CREATE de inscripción falló: %v", err)
	}

	if err := queries.DeleteInscripcion(ctx, db.DeleteInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       alumno.ID,
	}); err != nil {
		t.Fatalf("DELETE de inscripción falló: %v", err)
	}

	inscriptos, err := queries.ListInscriptos(ctx, horario.ID)
	if err != nil {
		t.Fatalf("LIST después del DELETE falló: %v", err)
	}
	for _, i := range inscriptos {
		if i.AlumnoID == alumno.ID {
			t.Fatal("el alumno no debió seguir inscripto después del DELETE")
		}
	}
}

func TestInscripcionDuplicada(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "lunes", horaDePrueba(8, 0))
	alumno := createAlumnoDePrueba(t)

	params := db.CreateInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       alumno.ID,
	}
	if _, err := queries.CreateInscripcion(ctx, params); err != nil {
		t.Fatalf("primera inscripción falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteInscripcion(ctx, db.DeleteInscripcionParams{
			ClaseHorarioID: horario.ID,
			AlumnoID:       alumno.ID,
		})
	})

	_, err := queries.CreateInscripcion(ctx, params)
	if err == nil {
		t.Fatal("no se debió poder inscribir al mismo alumno dos veces en el mismo horario")
	}
	if !esCodigoPostgres(err, "23505") {
		t.Fatalf("se esperaba unique violation (23505), se obtuvo: %v", err)
	}
}

func TestInscripcionConAlumnoInexistente(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "martes", horaDePrueba(8, 30))

	_, err := queries.CreateInscripcion(ctx, db.CreateInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       -1,
	})
	if err == nil {
		t.Fatal("CREATE debió fallar por FK de alumno_id")
	}
	if !esCodigoPostgres(err, "23503") {
		t.Fatalf("se esperaba violación de FK (23503), se obtuvo: %v", err)
	}
}

func TestDeleteAlumnoCascadaInscripcion(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "miercoles", horaDePrueba(16, 0))
	alumno := createAlumnoDePrueba(t)

	if _, err := queries.CreateInscripcion(ctx, db.CreateInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       alumno.ID,
	}); err != nil {
		t.Fatalf("CREATE de inscripción falló: %v", err)
	}

	if err := queries.DeleteUsuario(ctx, alumno.ID); err != nil {
		t.Fatalf("DELETE del alumno falló: %v", err)
	}

	inscriptos, err := queries.ListInscriptos(ctx, horario.ID)
	if err != nil {
		t.Fatalf("LIST de inscriptos falló: %v", err)
	}
	for _, i := range inscriptos {
		if i.AlumnoID == alumno.ID {
			t.Fatal("la inscripción debió eliminarse en cascada al borrar el alumno")
		}
	}

	_, err = queries.GetAlumnoCompleto(ctx, alumno.ID)
	if err == nil {
		t.Fatal("el alumno no debió existir")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado: %v", err)
	}
}

func TestDeleteClaseCascadaInscripcion(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	clase := createClaseDePrueba(t, profesor.ID)
	horario := createHorarioDePrueba(t, clase.ID, "jueves", horaDePrueba(20, 0))
	alumno := createAlumnoDePrueba(t)

	if _, err := queries.CreateInscripcion(ctx, db.CreateInscripcionParams{
		ClaseHorarioID: horario.ID,
		AlumnoID:       alumno.ID,
	}); err != nil {
		t.Fatalf("CREATE de inscripción falló: %v", err)
	}

	if err := queries.DeleteClase(ctx, clase.ID); err != nil {
		t.Fatalf("DELETE de clase falló: %v", err)
	}

	inscriptos, err := queries.ListInscriptos(ctx, horario.ID)
	if err != nil {
		t.Fatalf("LIST de inscriptos falló: %v", err)
	}
	if len(inscriptos) != 0 {
		t.Fatal("borrar la clase debió eliminar horarios e inscripciones en cascada")
	}
}
