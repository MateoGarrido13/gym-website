package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"
)

func createUsuarioDePrueba(t *testing.T) db.CreateUsuarioRow {
	t.Helper()

	params := db.CreateUsuarioParams{ // creamos un struct con los parametros para el createUsuario
		Email:      fmt.Sprintf("test-%d@example.com", time.Now().UnixNano()),
		Contrasena: "secret123",
		Nombre:     "Mateo",
		Apellido:   "Compilado",
		Telefono:   sql.NullString{String: "123456789", Valid: true},
		Rol:        "alumno",
	}

	usuario, err := queries.CreateUsuario(ctx, params) // pasamos el contexto e intentamos una insercion con los parametros
	if err != nil {
		t.Fatalf("error al crear usuario de prueba: %v", err)
	}

	t.Cleanup(func() {
		_ = queries.DeleteUsuario(ctx, usuario.ID)
	})

	return usuario
}

func TestCreateUsuario(t *testing.T) {
	params := db.CreateUsuarioParams{ // creamos un struct con los parametros para el createUsuario
		Email:      fmt.Sprintf("create-%d@example.com", time.Now().UnixNano()),
		Contrasena: "secret123",
		Nombre:     "Mateo",
		Apellido:   "Compilado",
		Telefono:   sql.NullString{String: "123456789", Valid: true},
		Rol:        "alumno",
	}

	usuario, err := queries.CreateUsuario(ctx, params)
	if err != nil {
		t.Fatalf("CREATE falló: %v", err)
	}

	if usuario.ID == 0 {
		t.Fatal("CREATE no devolvió un ID válido")
	}
	if !usuario.CreatedAt.Valid {
		t.Error("created_at debería estar definido")
	}
	t.Cleanup(func() {
		_ = queries.DeleteUsuario(ctx, usuario.ID)
	})

	/* PODRIAMOS HACER UNA PRUEBA MAS EXTENSA DE CADA ATRIBUTO SI SE VERIFICA QUE LOS VALORES SON CORRECTOS
	if usuario.Email != params.Email {
		t.Errorf("email = %q, se esperaba %q", usuario.Email, params.Email)
	}
	if usuario.Nombre != params.Nombre || usuario.Apellido != params.Apellido {
		t.Errorf("nombre = %q %q, se esperaba %q %q", usuario.Nombre, usuario.Apellido, params.Nombre, params.Apellido)
	}
	if usuario.Rol != params.Rol {
		t.Errorf("rol = %q, se esperaba %q", usuario.Rol, params.Rol)
	} */

}

func TestGetUsuario(t *testing.T) {
	creado := createUsuarioDePrueba(t)

	usuario, err := queries.GetUsuario(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ falló: %v", err)
	}

	if usuario.ID != creado.ID {
		t.Errorf("id = %d, se esperaba %d", usuario.ID, creado.ID)
	}
	if usuario.Email != creado.Email {
		t.Errorf("email = %q, se esperaba %q", usuario.Email, creado.Email)
	}
	if usuario.Nombre != creado.Nombre || usuario.Apellido != creado.Apellido {
		t.Errorf("nombre = %q %q, se esperaba %q %q", usuario.Nombre, usuario.Apellido, creado.Nombre, creado.Apellido)
	}
	if usuario.Contrasena != "secret123" {
		t.Errorf("contrasena = %q, se esperaba %q", usuario.Contrasena, "secret123")
	}
}

func TestUpdateUsuario(t *testing.T) {
	creado := createUsuarioDePrueba(t)

	updateParams := db.UpdateUsuarioParams{
		ID:         creado.ID,
		Email:      fmt.Sprintf("actualizado-%d@example.com", time.Now().UnixNano()),
		Contrasena: "nuevaClave456",
		Nombre:     "Mateo",
		Apellido:   "Actualizado",
		Telefono:   sql.NullString{String: "987654321", Valid: true},
	}

	if err := queries.UpdateUsuario(ctx, updateParams); err != nil {
		t.Fatalf("UPDATE falló: %v", err)
	}

	usuario, err := queries.GetUsuario(ctx, creado.ID)
	if err != nil {
		t.Fatalf("READ después del UPDATE falló: %v", err)
	}

	if usuario.Email != updateParams.Email {
		t.Errorf("email = %q, se esperaba %q", usuario.Email, updateParams.Email)
	}
	if usuario.Apellido != updateParams.Apellido {
		t.Errorf("apellido = %q, se esperaba %q", usuario.Apellido, updateParams.Apellido)
	}
	if usuario.Contrasena != updateParams.Contrasena {
		t.Errorf("contrasena = %q, se esperaba %q", usuario.Contrasena, updateParams.Contrasena)
	}
	if usuario.Telefono.String != updateParams.Telefono.String {
		t.Errorf("telefono = %q, se esperaba %q", usuario.Telefono.String, updateParams.Telefono.String)
	}
}

func TestDeleteUsuario(t *testing.T) {
	creado := createUsuarioDePrueba(t)

	if err := queries.DeleteUsuario(ctx, creado.ID); err != nil {
		t.Fatalf("DELETE falló: %v", err)
	}

	_, err := queries.GetUsuario(ctx, creado.ID)
	if err == nil {
		t.Fatal("el usuario no debió existir después del DELETE")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado al verificar el DELETE: %v", err)
	}
}

func TestCreateUsuarioProfesor(t *testing.T) {
	params := db.CreateUsuarioParams{
		Email:      fmt.Sprintf("profesor-user-%d@example.com", time.Now().UnixNano()),
		Contrasena: "secret123",
		Nombre:     "Ana",
		Apellido:   "Coach",
		Telefono:   sql.NullString{Valid: false},
		Rol:        "profesor",
	}

	usuario, err := queries.CreateUsuario(ctx, params)
	if err != nil {
		t.Fatalf("CREATE de usuario profesor falló: %v", err)
	}
	t.Cleanup(func() {
		_ = queries.DeleteUsuario(ctx, usuario.ID)
	})

	if usuario.Rol != "profesor" {
		t.Errorf("rol = %q, se esperaba %q", usuario.Rol, "profesor")
	}
	if usuario.Email != params.Email {
		t.Errorf("email = %q, se esperaba %q", usuario.Email, params.Email)
	}
	if usuario.Telefono.Valid {
		t.Error("telefono debía ser NULL")
	}
}

func TestCreateUsuarioEmailDuplicado(t *testing.T) {
	creado := createUsuarioDePrueba(t)

	_, err := queries.CreateUsuario(ctx, db.CreateUsuarioParams{
		Email:      creado.Email,
		Contrasena: "otraClave",
		Nombre:     "Otro",
		Apellido:   "Usuario",
		Rol:        "alumno",
	})
	if err == nil {
		t.Fatal("CREATE debió fallar por email duplicado")
	}
	if !esCodigoPostgres(err, "23505") {
		t.Fatalf("se esperaba unique violation (23505), se obtuvo: %v", err)
	}
}

func TestCreateUsuarioRolInvalido(t *testing.T) {
	_, err := queries.CreateUsuario(ctx, db.CreateUsuarioParams{
		Email:      fmt.Sprintf("rol-invalido-%d@example.com", time.Now().UnixNano()),
		Contrasena: "secret123",
		Nombre:     "Rol",
		Apellido:   "Invalido",
		Rol:        "admin",
	})
	if err == nil {
		t.Fatal("CREATE debió fallar por rol inválido")
	}
	if !esCodigoPostgres(err, "22P02") {
		t.Fatalf("se esperaba invalid_text_representation (22P02), se obtuvo: %v", err)
	}
}

func TestGetUsuarioInexistente(t *testing.T) {
	_, err := queries.GetUsuario(ctx, -1)
	if err == nil {
		t.Fatal("GetUsuario debió fallar con un id inexistente")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error inesperado: %v", err)
	}
}
