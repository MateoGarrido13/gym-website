package test

import (
	"PRACTICO_DOS/internal/handler"
	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/internal/web"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func nuevoRouterDePrueba() http.Handler {
	return handler.NewRouter(
		service.NewEjercicioService(queries), service.NewAlumnoService(queries),
		service.NewProfesorService(queries), service.NewRutinaService(queries),
		service.NewRutinaEjercicioService(queries), service.NewClaseService(queries),
		service.NewClaseHorarioService(queries), service.NewInscripcionService(queries),
		service.NewAuthService(queries), "../static",
	)
}

func postLogin(router http.Handler, email, clave, rol string) *httptest.ResponseRecorder {
	form := url.Values{"email": {email}, "contrasena": {clave}, "rol": {rol}}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func cookieDeSesion(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == "session" {
			return c
		}
	}
	return nil
}

func TestLoginPageMuestraLosDosAccesos(t *testing.T) {
	rec := httptest.NewRecorder()
	nuevoRouterDePrueba().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/login", nil))

	body := rec.Body.String()
	if rec.Code != http.StatusOK ||
		!strings.Contains(body, "Entrar como profesor") || !strings.Contains(body, "Entrar como alumno") {
		t.Fatalf("status %d, faltan los accesos por rol", rec.Code)
	}
}

func TestLoginPageConRolMuestraElFormulario(t *testing.T) {
	rec := httptest.NewRecorder()
	nuevoRouterDePrueba().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/login?rol=alumno", nil))

	body := rec.Body.String()
	if !strings.Contains(body, "Entrar como alumno") || !strings.Contains(body, `name="contrasena"`) ||
		!strings.Contains(body, `name="rol" value="alumno"`) {
		t.Fatalf("no se mostró el formulario del alumno:\n%s", body)
	}
}

func TestLoginWebCorrectoRedirigeSegunRolYSeteaCookie(t *testing.T) {
	router := nuevoRouterDePrueba()
	casos := []struct{ rol, destino string }{
		{service.RolProfesor, "/alumnos"},
		{service.RolAlumno, "/mi-rutina"},
	}
	for _, c := range casos {
		t.Run(c.rol, func(t *testing.T) {
			_, email := crearUsuarioConClave(t, c.rol, "secret123")
			rec := postLogin(router, email, "secret123", c.rol)

			if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != c.destino {
				t.Fatalf("status %d, Location %q; se esperaba 303 a %q", rec.Code, rec.Header().Get("Location"), c.destino)
			}
			cookie := cookieDeSesion(rec)
			if cookie == nil || cookie.Value == "" {
				t.Fatal("no se seteó la cookie de sesión")
			}
			if !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" {
				t.Errorf("flags de cookie incorrectos: %+v", cookie)
			}
		})
	}
}

func TestLoginWebFalla(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, email := crearUsuarioConClave(t, service.RolAlumno, "secret123")

	casos := []struct {
		nombre, email, clave, rol string
		status                    int
	}{
		{"contraseña incorrecta", email, "otra", service.RolAlumno, http.StatusUnauthorized},
		{"rol cruzado", email, "secret123", service.RolProfesor, http.StatusUnauthorized},
		{"campos vacíos", "", "", service.RolAlumno, http.StatusBadRequest},
		{"sin rol", email, "secret123", "", http.StatusBadRequest},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := postLogin(router, c.email, c.clave, c.rol)
			if rec.Code != c.status {
				t.Fatalf("status %d, se esperaba %d", rec.Code, c.status)
			}
			if cookieDeSesion(rec) != nil {
				t.Error("no debería setearse cookie si el login falla")
			}
			if strings.Contains(rec.Body.String(), "secret123") {
				t.Error("la página no debe devolver la contraseña")
			}
		})
	}
}

func TestLogoutWebCierraLaSesion(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, email := crearUsuarioConClave(t, service.RolAlumno, "secret123")
	cookie := cookieDeSesion(postLogin(router, email, "secret123", service.RolAlumno))

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/login" {
		t.Fatalf("status %d, Location %q", rec.Code, rec.Header().Get("Location"))
	}
	if c := cookieDeSesion(rec); c == nil || c.MaxAge >= 0 {
		t.Error("la cookie debería quedar borrada (MaxAge < 0)")
	}
	if _, err := service.NewAuthService(queries).UsuarioDeSesion(ctx, cookie.Value); err == nil {
		t.Error("el token debería haber dejado de valer en el servidor")
	}
}

func TestRequireRole(t *testing.T) {
	pages := web.NewHandler(nil, nil, nil, nil, nil, nil, service.NewAuthService(queries))
	router := nuevoRouterDePrueba()
	solo := pages.RequireRole(service.RolProfesor, func(w http.ResponseWriter, r *http.Request) {
		user, ok := web.UserFrom(r.Context())
		if !ok || user.Rol != service.RolProfesor {
			t.Errorf("el usuario no llegó por el contexto: %+v ok=%v", user, ok)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	_, emailProfe := crearUsuarioConClave(t, service.RolProfesor, "secret123")
	_, emailAlumno := crearUsuarioConClave(t, service.RolAlumno, "secret123")
	profe := cookieDeSesion(postLogin(router, emailProfe, "secret123", service.RolProfesor))
	alumno := cookieDeSesion(postLogin(router, emailAlumno, "secret123", service.RolAlumno))

	pedir := func(c *http.Cookie, hx bool) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		if c != nil {
			req.AddCookie(c)
		}
		if hx {
			req.Header.Set("HX-Request", "true")
		}
		rec := httptest.NewRecorder()
		solo(rec, req)
		return rec
	}

	if rec := pedir(profe, false); rec.Code != http.StatusNoContent {
		t.Errorf("profesor: status %d, se esperaba 204", rec.Code)
	}
	if rec := pedir(nil, false); rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/login" {
		t.Errorf("sin sesión: status %d, Location %q", rec.Code, rec.Header().Get("Location"))
	}
	if rec := pedir(&http.Cookie{Name: "session", Value: "falsa"}, false); rec.Header().Get("Location") != "/login" {
		t.Errorf("token inválido: Location %q", rec.Header().Get("Location"))
	}
	if rec := pedir(alumno, false); rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/mi-rutina" {
		t.Errorf("alumno en ruta de profesor: status %d, Location %q", rec.Code, rec.Header().Get("Location"))
	}
	rec := pedir(nil, true)
	if rec.Code != http.StatusUnauthorized || rec.Header().Get("HX-Redirect") != "/login" {
		t.Errorf("HTMX sin sesión: status %d, HX-Redirect %q", rec.Code, rec.Header().Get("HX-Redirect"))
	}
}

func pedirConCookie(router http.Handler, metodo, ruta string, c *http.Cookie, cabeceras map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, ruta, nil)
	if c != nil {
		req.AddCookie(c)
	}
	for k, v := range cabeceras {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestHomeRedirigeSegunSesion(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, emailProfe := crearUsuarioConClave(t, service.RolProfesor, "secret123")
	_, emailAlumno := crearUsuarioConClave(t, service.RolAlumno, "secret123")
	profe := cookieDeSesion(postLogin(router, emailProfe, "secret123", service.RolProfesor))
	alumno := cookieDeSesion(postLogin(router, emailAlumno, "secret123", service.RolAlumno))

	casos := []struct {
		nombre  string
		cookie  *http.Cookie
		destino string
	}{
		{"sin sesión", nil, "/login"},
		{"profesor", profe, "/alumnos"},
		{"alumno", alumno, "/mi-rutina"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := pedirConCookie(router, http.MethodGet, "/", c.cookie, nil)
			if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != c.destino {
				t.Fatalf("status %d, Location %q; se esperaba 303 a %q", rec.Code, rec.Header().Get("Location"), c.destino)
			}
		})
	}

	if rec := pedirConCookie(router, http.MethodGet, "/no-existe", profe, nil); rec.Code != http.StatusNotFound {
		t.Errorf("ruta desconocida: status %d, se esperaba 404", rec.Code)
	}
}

func TestCatalogoSoloParaProfesores(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, emailProfe := crearUsuarioConClave(t, service.RolProfesor, "secret123")
	_, emailAlumno := crearUsuarioConClave(t, service.RolAlumno, "secret123")
	profe := cookieDeSesion(postLogin(router, emailProfe, "secret123", service.RolProfesor))
	alumno := cookieDeSesion(postLogin(router, emailAlumno, "secret123", service.RolAlumno))

	rutas := []struct{ metodo, ruta string }{
		{"GET", "/alumnos"}, {"POST", "/alumnos"}, {"GET", "/alumnos/1/editar"},
		{"POST", "/alumnos/1"}, {"DELETE", "/alumnos/1"},
		{"GET", "/profesores"}, {"POST", "/profesores"}, {"DELETE", "/profesores/1"},
		{"GET", "/ejercicios"}, {"POST", "/ejercicios"}, {"DELETE", "/ejercicios/1"},
		{"GET", "/rutinas"}, {"POST", "/rutinas"}, {"DELETE", "/rutinas/1"},
		{"GET", "/clases"}, {"POST", "/clases"}, {"DELETE", "/clases/1"},
	}
	for _, r := range rutas {
		t.Run(r.metodo+" "+r.ruta, func(t *testing.T) {
			if rec := pedirConCookie(router, r.metodo, r.ruta, nil, nil); rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/login" {
				t.Errorf("sin sesión: status %d, Location %q", rec.Code, rec.Header().Get("Location"))
			}
			if rec := pedirConCookie(router, r.metodo, r.ruta, alumno, nil); rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/mi-rutina" {
				t.Errorf("alumno: status %d, Location %q", rec.Code, rec.Header().Get("Location"))
			}
		})
	}

	// HTMX sin sesión no debe recibir un redirect que meta el login en un fragmento.
	rec := pedirConCookie(router, "DELETE", "/alumnos/1", nil, map[string]string{"HX-Request": "true"})
	if rec.Code != http.StatusUnauthorized || rec.Header().Get("HX-Redirect") != "/login" {
		t.Errorf("HTMX: status %d, HX-Redirect %q", rec.Code, rec.Header().Get("HX-Redirect"))
	}

	rec = pedirConCookie(router, "GET", "/alumnos", profe, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("profesor: status %d, se esperaba 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Pablo Profe") || !strings.Contains(body, `action="/logout"`) {
		t.Error("la página del profesor debería mostrar su nombre y el botón Salir")
	}
}

func TestPostSinSesionNoCreaDatos(t *testing.T) {
	router := nuevoRouterDePrueba()
	email := "intruso-sin-sesion@example.com"
	form := url.Values{"nombre": {"X"}, "apellido": {"Y"}, "email": {email}, "contrasena": {"secret123"}, "fecha_inscripcion": {"2026-01-01"}}
	req := httptest.NewRequest(http.MethodPost, "/alumnos", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(httptest.NewRecorder(), req)

	if _, err := queries.GetUsuarioByEmail(ctx, email); err == nil {
		t.Error("un POST sin sesión no debería crear el alumno")
	}
}
