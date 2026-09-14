package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
)


func horaDePrueba(hora, minuto int) time.Time {
	return time.Date(1970, 1, 1, hora, minuto, 0, 0, time.UTC)
}

func mismaHora(a, b time.Time) bool {
	return a.Hour() == b.Hour() && a.Minute() == b.Minute() && a.Second() == b.Second()
}

func esCodigoPostgres(err error, codigo string) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && string(pqErr.Code) == codigo
}

// helpers para los tests
func createProfesorDePrueba(t *testing.T) db.GetProfesorCompletoRow {
	t.Helper()

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	params := db.CreateProfesorParams{
		Email:        fmt.Sprintf("profesor-%s@example.com", suffix),
		Contrasena:   "secret123",
		Nombre:       "Laura",
		Apellido:     fmt.Sprintf("Profesora-%s", suffix),
		Telefono:     sql.NullString{String: "111111111", Valid: true},
		Especialidad: sql.NullString{String: "fuerza", Valid: true},
	}

	profesor, err := queries.CreateProfesor(ctx, params)
	if err != nil {
		t.Fatalf("error al crear profesor de prueba: %v", err)
	}

	t.Cleanup(func() {
		_ = queries.DeleteUsuario(ctx, profesor.ID)
	})

	return profesor
}

func createRutinaDePrueba(t *testing.T, profesorID int32) db.Rutina {
	t.Helper()

	rutina, err := queries.CreateRutina(ctx, db.CreateRutinaParams{
		Nombre:          fmt.Sprintf("rutina-%d", time.Now().UnixNano()),
		DuracionSemanas: 8,
		ProfesorID:      profesorID,
	})
	if err != nil {
		t.Fatalf("error al crear rutina de prueba: %v", err)
	}

	t.Cleanup(func() {
		_ = queries.DeleteRutina(ctx, rutina.ID)
	})

	return rutina
}

func createClaseDePrueba(t *testing.T, profesorID int32) db.Clase {
	t.Helper()

	clase, err := queries.CreateClase(ctx, db.CreateClaseParams{
		Nombre:     fmt.Sprintf("spinning-%d", time.Now().UnixNano()),
		MaxAlumnos: 12,
		ProfesorID: profesorID,
	})
	if err != nil {
		t.Fatalf("error al crear clase de prueba: %v", err)
	}

	t.Cleanup(func() {
		_ = queries.DeleteClase(ctx, clase.ID)
	})

	return clase
}

func createHorarioDePrueba(t *testing.T, claseID int32, dia string, hora time.Time) db.ClaseHorario {
	t.Helper()

	horario, err := queries.CreateClaseHorario(ctx, db.CreateClaseHorarioParams{
		ClaseID:    claseID,
		DiaSemana:  dia,
		HoraInicio: hora,
	})
	if err != nil {
		t.Fatalf("error al crear horario de prueba: %v", err)
	}

	t.Cleanup(func() {
		_ = queries.DeleteClaseHorario(ctx, horario.ID)
	})

	return horario
}
