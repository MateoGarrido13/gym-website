package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"PRACTICO_DOS/internal/service"
)

func TestCreateRutinaValidaRenderizaFormularioVacioConExito(t *testing.T) {
	router := nuevoRouterDePrueba()
	_, email := crearUsuarioConClave(t, service.RolProfesor, "secret123")
	cookie := cookieDeSesion(postLogin(router, email, "secret123", service.RolProfesor))
	if cookie == nil {
		t.Fatal("el login no devolvió cookie de sesión")
	}

	profesor := createProfesorDePrueba(t)
	ejercicio := createEjercicioDePrueba(t)
	const nombre = "Rutina de prueba"
	form := url.Values{
		"nombre":           {nombre},
		"duracion_semanas": {"8"},
		"profesor_id":      {fmt.Sprintf("%d", profesor.ID)},
		"ejercicio_id":     {fmt.Sprintf("%d", ejercicio.ID)},
	}

	cases := []struct {
		name string
		htmx bool
	}{
		{name: "sin HTMX", htmx: false},
		{name: "con HTMX", htmx: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/rutinas", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if tc.htmx {
				req.Header.Set("HX-Request", "true")
			}
			req.AddCookie(cookie)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status %d, se esperaba 200: %s", rec.Code, rec.Body.String())
			}
			body := rec.Body.String()
			if !strings.Contains(body, "Registro creado correctamente") {
				t.Error("el cuerpo no contiene el mensaje de éxito")
			}
			if strings.Contains(body, `value="`+nombre+`"`) {
				t.Error("el campo nombre no debería conservar el valor anterior")
			}
			if !strings.Contains(body, nombre) {
				t.Error("la lista no muestra la rutina creada")
			}
			if !tc.htmx {
				return
			}
			if got := rec.Header().Get("HX-Retarget"); got != "body" {
				t.Errorf("HX-Retarget = %q, se esperaba body", got)
			}
			if got := rec.Header().Get("HX-Reswap"); got != "innerHTML" {
				t.Errorf("HX-Reswap = %q, se esperaba innerHTML", got)
			}
			if got := rec.Header().Get("HX-Reselect"); got != "" {
				t.Errorf("HX-Reselect = %q, no debe enviarse: seleccionar body vacía la página", got)
			}
		})
	}
}
