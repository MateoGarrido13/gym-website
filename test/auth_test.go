package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"PRACTICO_DOS/internal/service"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
	"time"
)

// crearUsuarioConClave crea un alumno, profesor o administrador por los servicios (que
// hashean la contraseña) y devuelve su email.
func crearUsuarioConClave(t *testing.T, rol, clave string) (int32, string) {
	t.Helper()

	email := fmt.Sprintf("auth-%s-%d@example.com", rol, time.Now().UnixNano())
	var id int32
	switch rol {
	case service.RolAdmin:
		// No hay servicio de alta de administradores: se inserta el usuario directo.
		hash, err := bcrypt.GenerateFromPassword([]byte(clave), bcrypt.MinCost)
		if err != nil {
			t.Fatalf("no se pudo hashear la clave: %v", err)
		}
		u, err := queries.CreateUsuario(ctx, db.CreateUsuarioParams{
			Email: email, Contrasena: string(hash), Nombre: "Adriana", Apellido: "Admin", Rol: service.RolAdmin,
		})
		if err != nil {
			t.Fatalf("no se pudo crear el administrador: %v", err)
		}
		id = u.ID
	case service.RolProfesor:
		p, err := service.NewProfesorService(queries).Create(ctx, service.ProfesorCreateInput{
			Email: email, Contrasena: clave, Nombre: "Pablo", Apellido: "Profe",
		})
		if err != nil {
			t.Fatalf("no se pudo crear el profesor: %v", err)
		}
		id = p.ID
	default:
		a, err := service.NewAlumnoService(queries).Create(ctx, service.AlumnoCreateInput{
			Email: email, Contrasena: clave, Nombre: "Ana", Apellido: "Alumna",
			FechaInscripcion: "2026-01-01",
		})
		if err != nil {
			t.Fatalf("no se pudo crear el alumno: %v", err)
		}
		id = a.ID
	}
	t.Cleanup(func() { _ = queries.DeleteUsuario(ctx, id) })
	return id, email
}

func TestAltaGuardaContrasenaHasheada(t *testing.T) {
	id, _ := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	usuario, err := queries.GetUsuario(ctx, id)
	if err != nil {
		t.Fatalf("GetUsuario falló: %v", err)
	}
	if usuario.Contrasena == "secret123" {
		t.Fatal("la contraseña se guardó en texto plano")
	}
	if !strings.HasPrefix(usuario.Contrasena, "$2") {
		t.Errorf("la contraseña no parece un hash bcrypt: %q", usuario.Contrasena)
	}
}

func TestAltaRechazaContrasenaLarga(t *testing.T) {
	_, err := service.NewAlumnoService(queries).Create(ctx, service.AlumnoCreateInput{
		Email: "larga@example.com", Contrasena: strings.Repeat("a", 73), Nombre: "A", Apellido: "B",
		FechaInscripcion: "2026-01-01",
	})
	if !errors.Is(err, service.ErrInvalid) {
		t.Fatalf("error = %v, se esperaba ErrInvalid", err)
	}
}

func TestLoginCorrecto(t *testing.T) {
	auth := service.NewAuthService(queries)
	id, email := crearUsuarioConClave(t, service.RolProfesor, "secret123")

	sesion, err := auth.Login(ctx, email, "secret123", service.RolProfesor)
	if err != nil {
		t.Fatalf("Login falló: %v", err)
	}
	if sesion.Token == "" || sesion.Usuario.ID != id || sesion.Usuario.Rol != service.RolProfesor {
		t.Fatalf("sesión inesperada: %+v", sesion)
	}
	if d := time.Until(sesion.Expira); d < 23*time.Hour || d > 24*time.Hour {
		t.Errorf("la sesión vence en %v, se esperaban 24h", d)
	}

	usuario, err := auth.UsuarioDeSesion(ctx, sesion.Token)
	if err != nil {
		t.Fatalf("UsuarioDeSesion falló: %v", err)
	}
	if usuario.ID != id {
		t.Errorf("usuario = %d, se esperaba %d", usuario.ID, id)
	}
}

func TestLoginGuardaHashDelTokenNoElToken(t *testing.T) {
	auth := service.NewAuthService(queries)
	_, email := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	sesion, err := auth.Login(ctx, email, "secret123", service.RolAlumno)
	if err != nil {
		t.Fatalf("Login falló: %v", err)
	}
	if _, err := queries.GetSesionUsuario(ctx, sesion.Token); err == nil {
		t.Fatal("el token en claro no debería existir en la tabla sesion")
	}
}

func TestLoginFallaConCredencialesIncorrectas(t *testing.T) {
	auth := service.NewAuthService(queries)
	_, email := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	casos := []struct {
		nombre, email, clave, rol string
	}{
		{"contraseña incorrecta", email, "otra", service.RolAlumno},
		{"rol cruzado", email, "secret123", service.RolProfesor},
		{"email inexistente", "nadie@example.com", "secret123", service.RolAlumno},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			_, err := auth.Login(ctx, c.email, c.clave, c.rol)
			if !errors.Is(err, service.ErrUnauthorized) {
				t.Fatalf("error = %v, se esperaba ErrUnauthorized", err)
			}
		})
	}
}

func TestLoginValidaEntrada(t *testing.T) {
	auth := service.NewAuthService(queries)

	for _, c := range []struct{ nombre, email, clave, rol string }{
		{"rol inválido", "a@example.com", "x", "admin"},
		{"email vacío", "", "x", service.RolAlumno},
		{"contraseña vacía", "a@example.com", "", service.RolAlumno},
	} {
		t.Run(c.nombre, func(t *testing.T) {
			if _, err := auth.Login(ctx, c.email, c.clave, c.rol); !errors.Is(err, service.ErrInvalid) {
				t.Fatalf("error = %v, se esperaba ErrInvalid", err)
			}
		})
	}
}

func TestSesionInexistenteOVacia(t *testing.T) {
	auth := service.NewAuthService(queries)
	for _, token := range []string{"", "no-existe"} {
		if _, err := auth.UsuarioDeSesion(ctx, token); !errors.Is(err, service.ErrUnauthorized) {
			t.Errorf("token %q: error = %v, se esperaba ErrUnauthorized", token, err)
		}
	}
}

func TestSesionExpiradaNoVale(t *testing.T) {
	auth := service.NewAuthService(queries)
	id, _ := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	// Se inserta directo una sesión que ya venció (el servicio siempre crea
	// sesiones a 24h). El id es un hash cualquiera: acá no importa el token.
	vencida := db.CreateSesionParams{ID: strings.Repeat("a", 64), UsuarioID: id, ExpiraEn: time.Now().Add(-time.Minute)}
	if err := queries.CreateSesion(ctx, vencida); err != nil {
		t.Fatalf("CreateSesion falló: %v", err)
	}
	if _, err := queries.GetSesionUsuario(ctx, vencida.ID); err == nil {
		t.Fatal("una sesión vencida no debería devolver usuario")
	}

	if err := auth.LimpiarExpiradas(ctx); err != nil {
		t.Fatalf("LimpiarExpiradas falló: %v", err)
	}
	var n int
	if err := dbConn.QueryRowContext(ctx, "SELECT COUNT(*) FROM sesion WHERE id = $1", vencida.ID).Scan(&n); err != nil || n != 0 {
		t.Errorf("la sesión vencida debería haberse borrado (n=%d, err=%v)", n, err)
	}
}

func TestLogoutCierraLaSesion(t *testing.T) {
	auth := service.NewAuthService(queries)
	_, email := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	sesion, err := auth.Login(ctx, email, "secret123", service.RolAlumno)
	if err != nil {
		t.Fatalf("Login falló: %v", err)
	}
	if err := auth.Logout(ctx, sesion.Token); err != nil {
		t.Fatalf("Logout falló: %v", err)
	}
	if _, err := auth.UsuarioDeSesion(ctx, sesion.Token); !errors.Is(err, service.ErrUnauthorized) {
		t.Fatalf("error = %v, se esperaba ErrUnauthorized tras el logout", err)
	}
	if err := auth.Logout(ctx, sesion.Token); err != nil {
		t.Errorf("Logout de una sesión ya cerrada no debería fallar: %v", err)
	}
}

func TestBorrarUsuarioCierraSusSesiones(t *testing.T) {
	auth := service.NewAuthService(queries)
	id, email := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	sesion, err := auth.Login(ctx, email, "secret123", service.RolAlumno)
	if err != nil {
		t.Fatalf("Login falló: %v", err)
	}
	if err := queries.DeleteUsuario(ctx, id); err != nil {
		t.Fatalf("DeleteUsuario falló: %v", err)
	}
	if _, err := auth.UsuarioDeSesion(ctx, sesion.Token); !errors.Is(err, service.ErrUnauthorized) {
		t.Fatalf("error = %v, se esperaba ErrUnauthorized", err)
	}
}
