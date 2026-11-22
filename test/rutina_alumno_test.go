package test

import (
	db "PRACTICO_DOS/db/sqlc"
	"PRACTICO_DOS/internal/service"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

// alumnoConRutina arma un profesor, una rutina con dos ejercicios (cargados
// a propósito en orden inverso al de "orden") y un alumno con esa rutina.
// Devuelve el email del alumno y la rutina.
func alumnoConRutina(t *testing.T) (string, db.Rutina) {
	t.Helper()
	// El orden de creación importa: los Cleanup corren al revés (alumno,
	// rutina, profesor, ejercicios), respetando las claves foráneas.
	ejercicioB := createEjercicioDePrueba(t)
	ejercicioA := createEjercicioDePrueba(t)
	profesor := createProfesorDePrueba(t)
	rutina := createRutinaDePrueba(t, profesor.ID)

	for _, item := range []db.AddEjercicioARutinaParams{
		{RutinaID: rutina.ID, EjercicioID: ejercicioB.ID, Orden: 2, Series: sql.NullInt32{Int32: 3, Valid: true}},
		{RutinaID: rutina.ID, EjercicioID: ejercicioA.ID, Orden: 1,
			Series: sql.NullInt32{Int32: 4, Valid: true}, Repeticiones: sql.NullInt32{Int32: 12, Valid: true},
			DescansoSegundos: sql.NullInt32{Int32: 60, Valid: true}},
	} {
		if err := queries.AddEjercicioARutina(ctx, item); err != nil {
			t.Fatalf("no se pudo agregar el ejercicio a la rutina: %v", err)
		}
	}

	email := fmt.Sprintf("alumno-rutina-%d@example.com", time.Now().UnixNano())
	alumno, err := service.NewAlumnoService(queries).Create(ctx, service.AlumnoCreateInput{
		Email: email, Contrasena: "secret123", Nombre: "Rita", Apellido: "Rutina",
		FechaInscripcion: "2026-01-01", RutinaID: &rutina.ID,
	})
	if err != nil {
		t.Fatalf("no se pudo crear el alumno: %v", err)
	}
	t.Cleanup(func() { _ = queries.DeleteUsuario(ctx, alumno.ID) })
	return email, rutina
}

func TestRutinaDeAlumno(t *testing.T) {
	email, rutina := alumnoConRutina(t)
	usuario, err := queries.GetUsuarioByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetUsuarioByEmail falló: %v", err)
	}

	got, err := service.NewRutinaService(queries).DeAlumno(ctx, usuario.ID)
	if err != nil {
		t.Fatalf("DeAlumno falló: %v", err)
	}
	if got.ID != rutina.ID || got.Nombre != rutina.Nombre || got.DuracionSemanas != 8 {
		t.Errorf("rutina inesperada: %+v", got)
	}
	if !strings.Contains(got.Profesor, "Laura") {
		t.Errorf("profesor = %q, se esperaba el nombre del profesor", got.Profesor)
	}
	if len(got.Ejercicios) != 2 {
		t.Fatalf("ejercicios = %d, se esperaban 2", len(got.Ejercicios))
	}
	primero, segundo := got.Ejercicios[0], got.Ejercicios[1]
	if primero.Orden != 1 || segundo.Orden != 2 {
		t.Errorf("orden = %d, %d; se esperaba 1, 2", primero.Orden, segundo.Orden)
	}
	if primero.Series == nil || *primero.Series != 4 || *primero.Repeticiones != 12 || *primero.DescansoSegundos != 60 {
		t.Errorf("detalle del primer ejercicio incorrecto: %+v", primero)
	}
	if segundo.Repeticiones != nil || segundo.DescansoSegundos != nil {
		t.Errorf("los datos no cargados deberían ser nil: %+v", segundo)
	}
}

func TestRutinaDeAlumnoSinRutina(t *testing.T) {
	id, _ := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	_, err := service.NewRutinaService(queries).DeAlumno(ctx, id)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("error = %v, se esperaba ErrNotFound", err)
	}
}

func TestMiRutinaWeb(t *testing.T) {
	router := nuevoRouterDePrueba()
	email, rutina := alumnoConRutina(t)
	cookie := cookieDeSesion(postLogin(router, email, "secret123", service.RolAlumno))

	rec := pedirConCookie(router, http.MethodGet, "/mi-rutina", cookie, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, se esperaba 200", rec.Code)
	}
	body := rec.Body.String()
	for _, esperado := range []string{rutina.Nombre, "8 semanas", "Rita Rutina", "4 series × 12 repeticiones", "60 s de descanso", "3 series", `action="/logout"`} {
		if !strings.Contains(body, esperado) {
			t.Errorf("la página no contiene %q", esperado)
		}
	}
	if strings.Index(body, "4 series") > strings.Index(body, "3 series") {
		t.Error("los ejercicios no están en el orden de la rutina")
	}
	if strings.Contains(body, "Alumnos") || strings.Contains(body, "/profesores") {
		t.Error("la vista del alumno no debe mostrar el menú de administración")
	}
}

func TestMiRutinaWebSinRutina(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, email := crearUsuarioConClave(t, service.RolAlumno, "secret123")
	cookie := cookieDeSesion(postLogin(router, email, "secret123", service.RolAlumno))

	rec := pedirConCookie(router, http.MethodGet, "/mi-rutina", cookie, nil)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Todavía no tienes una rutina asignada") {
		t.Fatalf("status %d; se esperaba el aviso de sin rutina", rec.Code)
	}
}

// Cada alumno ve solo su rutina: B, sin rutina, no ve la de A.
func TestMiRutinaWebNoMuestraLaRutinaDeOtro(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, rutinaDeA := alumnoConRutina(t)
	_, emailB := crearUsuarioConClave(t, service.RolAlumno, "secret123")
	cookieB := cookieDeSesion(postLogin(router, emailB, "secret123", service.RolAlumno))

	rec := pedirConCookie(router, http.MethodGet, "/mi-rutina", cookieB, nil)
	if strings.Contains(rec.Body.String(), rutinaDeA.Nombre) {
		t.Error("un alumno no debería ver la rutina de otro")
	}
}

func TestMiRutinaWebSoloParaAlumnos(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, emailProfe := crearUsuarioConClave(t, service.RolProfesor, "secret123")
	profe := cookieDeSesion(postLogin(router, emailProfe, "secret123", service.RolProfesor))

	if rec := pedirConCookie(router, http.MethodGet, "/mi-rutina", nil, nil); rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/login" {
		t.Errorf("sin sesión: status %d, Location %q", rec.Code, rec.Header().Get("Location"))
	}
	if rec := pedirConCookie(router, http.MethodGet, "/mi-rutina", profe, nil); rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/alumnos" {
		t.Errorf("profesor: status %d, Location %q", rec.Code, rec.Header().Get("Location"))
	}
}
