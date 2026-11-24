package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"
)

func createAlumnoDePrueba(t *testing.T) db.GetAlumnoCompletoRow {
	t.Helper()

	params := db.CreateAlumnoParams{
		Email:            fmt.Sprintf("alumno-%d@example.com", time.Now().UnixNano()),
		Contrasena:       "secret123",
		Nombre:           "Mateo",
		Apellido:         "Alumno",
		Telefono:         sql.NullString{String: "123456789", Valid: true},
		FechaInscripcion: time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC),
		FechaVto:         sql.NullTime{Time: time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		TipoPlan:         db.NullEnumTipoPlan{EnumTipoPlan: db.EnumTipoPlanMensual, Valid: true},
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
		TipoPlan:         db.NullEnumTipoPlan{EnumTipoPlan: db.EnumTipoPlanMensual, Valid: true},
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
	if alumno.TipoPlan.EnumTipoPlan != params.TipoPlan.EnumTipoPlan {
		t.Errorf("tipo_plan = %q, se esperaba %q", alumno.TipoPlan.EnumTipoPlan, params.TipoPlan.EnumTipoPlan)
	}
	if !mismaFecha(alumno.FechaInscripcion, fechaInscripcion) {
		t.Errorf("fecha_inscripcion = %s, se esperaba %s", alumno.FechaInscripcion.Format("2006-01-02"), fechaInscripcion.Format("2006-01-02"))
	}
}

func TestGetAlumno(t *testing.T) {
	creado := createAlumnoDePrueba(t)

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
	if alumno.TipoPlan.EnumTipoPlan != creado.TipoPlan.EnumTipoPlan {
		t.Errorf("tipo_plan = %q, se esperaba %q", alumno.TipoPlan.EnumTipoPlan, creado.TipoPlan.EnumTipoPlan)
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
	creado := createAlumnoDePrueba(t)

	nuevaFechaVto := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	updateParams := db.UpdateAlumnoDetallesParams{
		UsuarioID: creado.ID,
		FechaVto:  sql.NullTime{Time: nuevaFechaVto, Valid: true},
		TipoPlan:  db.NullEnumTipoPlan{EnumTipoPlan: db.EnumTipoPlanAnual, Valid: true},
		RutinaID:  sql.NullInt32{Valid: false},
	}

	if err := queries.UpdateAlumnoDetalles(ctx, updateParams); err != nil {
		t.Fatalf("UPDATE falló: %v", err)
	}

	alumno, err := queries.GetAlumnoCompleto(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ después del UPDATE falló: %v", err)
	}

	if alumno.TipoPlan.EnumTipoPlan != updateParams.TipoPlan.EnumTipoPlan {
		t.Errorf("tipo_plan = %q, se esperaba %q", alumno.TipoPlan.EnumTipoPlan, updateParams.TipoPlan.EnumTipoPlan)
	}
	if !alumno.FechaVto.Valid || !mismaFecha(alumno.FechaVto.Time, nuevaFechaVto) {
		t.Errorf("fecha_vto = %v, se esperaba %s", alumno.FechaVto, nuevaFechaVto.Format("2006-01-02"))
	}
	if alumno.Email != creado.Email {
		t.Errorf("el email del usuario padre no debió cambiar: %q", alumno.Email)
	}
}

func TestDeleteAlumno(t *testing.T) {
	creado := createAlumnoDePrueba(t)

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

func createAlumnoConNombre(t *testing.T, nombre, apellido, marca string) db.GetAlumnoCompletoRow {
	t.Helper()

	params := db.CreateAlumnoParams{
		Email:            fmt.Sprintf("orden-%s-%d@example.com", marca, time.Now().UnixNano()),
		Contrasena:       "secret123",
		Nombre:           nombre,
		Apellido:         apellido,
		Telefono:         sql.NullString{String: "123456789", Valid: true},
		FechaInscripcion: time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC),
		FechaVto:         sql.NullTime{Time: time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		TipoPlan:         db.NullEnumTipoPlan{EnumTipoPlan: db.EnumTipoPlanMensual, Valid: true},
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

func TestListAlumnosOrdenYPagina(t *testing.T) {
	anaMenor := createAlumnoConNombre(t, "Ana", "Perez", "menor")
	zoe := createAlumnoConNombre(t, "Zoe", "Perez", "zoe")
	anaMayor := createAlumnoConNombre(t, "Ana", "Perez", "mayor")
	if anaMenor.ID >= anaMayor.ID {
		t.Fatalf("id de Ana menor = %d, id de Ana mayor = %d", anaMenor.ID, anaMayor.ID)
	}

	const ventana int32 = 100
	var todos []db.ListAlumnosPaginaRow
	for offset := int32(0); ; offset += ventana {
		filas, err := queries.ListAlumnosPagina(ctx, db.ListAlumnosPaginaParams{
			Limit:  ventana,
			Offset: offset,
		})
		if err != nil {
			t.Fatalf("LIST paginado falló: %v", err)
		}
		todos = append(todos, filas...)
		if len(filas) < int(ventana) {
			break
		}
	}

	indice := map[int32]int{}
	for i, alumno := range todos {
		indice[alumno.ID] = i
	}
	iAnaMenor, okMenor := indice[anaMenor.ID]
	iAnaMayor, okMayor := indice[anaMayor.ID]
	iZoe, okZoe := indice[zoe.ID]
	if !okMenor || !okMayor || !okZoe {
		t.Fatal("algún alumno de prueba no apareció en el listado ordenado")
	}
	if iAnaMenor >= iAnaMayor || iAnaMayor >= iZoe {
		t.Fatalf("orden = ana menor %d, ana mayor %d, zoe %d", iAnaMenor, iAnaMayor, iZoe)
	}

	una, err := queries.ListAlumnosPagina(ctx, db.ListAlumnosPaginaParams{
		Limit:  1,
		Offset: int32(iAnaMayor),
	})
	if err != nil {
		t.Fatalf("ventana LIMIT 1 OFFSET %d falló: %v", iAnaMayor, err)
	}
	if len(una) != 1 || una[0].ID != anaMayor.ID {
		t.Fatalf("la ventana devolvió %v, se esperaba el alumno %d", idsDe(una), anaMayor.ID)
	}

	antes, err := queries.CountAlumnosAntesDe(ctx, anaMayor.ID)
	if err != nil {
		t.Fatalf("conteo de alumnos anteriores falló: %v", err)
	}
	if int(antes) != iAnaMayor {
		t.Errorf("alumnos anteriores = %d, se esperaba el índice %d", antes, iAnaMayor)
	}
}

func idsDe(filas []db.ListAlumnosPaginaRow) []int32 {
	ids := make([]int32, 0, len(filas))
	for _, fila := range filas {
		ids = append(ids, fila.ID)
	}
	return ids
}

func TestListAlumnos(t *testing.T) {
	creado := createAlumnoDePrueba(t)

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
			if a.TipoPlan.EnumTipoPlan != creado.TipoPlan.EnumTipoPlan {
				t.Errorf("tipo_plan = %q, se esperaba %q", a.TipoPlan.EnumTipoPlan, creado.TipoPlan.EnumTipoPlan)
			}
			break
		}
	}
	if !encontrado {
		t.Fatalf("el alumno creado (ID %d) no apareció en el listado", creado.ID)
	}
}

func TestAsignarRutinaAAlumno(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)
	alumno := createAlumnoDePrueba(t)

	if err := queries.UpdateAlumnoDetalles(ctx, db.UpdateAlumnoDetallesParams{
		UsuarioID: alumno.ID,
		FechaVto:  alumno.FechaVto,
		TipoPlan:  alumno.TipoPlan,
		RutinaID:  sql.NullInt32{Int32: rutina.ID, Valid: true},
	}); err != nil {
		t.Fatalf("asignar rutina falló: %v", err)
	}

	actualizado, err := queries.GetAlumnoCompleto(ctx, alumno.ID)
	if err != nil {
		t.Fatalf("READ después de asignar rutina falló: %v", err)
	}
	if !actualizado.RutinaID.Valid || actualizado.RutinaID.Int32 != rutina.ID {
		t.Errorf("rutina_id = %v, se esperaba %d", actualizado.RutinaID, rutina.ID)
	}
}

func TestAsignarRutinaInexistenteAAlumno(t *testing.T) {
	alumno := createAlumnoDePrueba(t)

	err := queries.UpdateAlumnoDetalles(ctx, db.UpdateAlumnoDetallesParams{
		UsuarioID: alumno.ID,
		FechaVto:  alumno.FechaVto,
		TipoPlan:  alumno.TipoPlan,
		RutinaID:  sql.NullInt32{Int32: -1, Valid: true},
	})
	if err == nil {
		t.Fatal("asignar una rutina inexistente debió fallar")
	}
	if !esCodigoPostgres(err, "23503") {
		t.Fatalf("se esperaba violación de FK (23503), se obtuvo: %v", err)
	}
}

func TestGetAlumnoInexistente(t *testing.T) {
	_, err := queries.GetAlumnoCompleto(ctx, -1)
	if err == nil {
		t.Fatal("GetAlumnoCompleto debió fallar con un id inexistente")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado: %v", err)
	}
}

func TestDeleteRutinaAsignadaAAlumno(t *testing.T) {
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)
	alumno := createAlumnoDePrueba(t)

	if err := queries.UpdateAlumnoDetalles(ctx, db.UpdateAlumnoDetallesParams{
		UsuarioID: alumno.ID,
		FechaVto:  alumno.FechaVto,
		TipoPlan:  alumno.TipoPlan,
		RutinaID:  sql.NullInt32{Int32: rutina.ID, Valid: true},
	}); err != nil {
		t.Fatalf("asignar rutina falló: %v", err)
	}

	err := queries.DeleteRutina(ctx, rutina.ID)
	if err == nil {
		t.Fatal("no se debió poder borrar una rutina asignada a un alumno")
	}
	if !esCodigoPostgres(err, "23503") {
		t.Fatalf("se esperaba violación de FK (23503), se obtuvo: %v", err)
	}
}
