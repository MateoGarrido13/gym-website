package web

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"PRACTICO_DOS/internal/service"
	"PRACTICO_DOS/views"
)

const sessionCookie = "session"

// homeFor es a donde va cada rol después de entrar.
func homeFor(rol string) string {
	if rol == service.RolProfesor {
		return "/alumnos"
	}
	return "/mi-rutina"
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

// currentUser resuelve el usuario de la cookie de sesión. ok=false significa
// "no hay sesión válida"; err solo viene si falló la base.
func (h *Handler) currentUser(r *http.Request) (user service.UsuarioSesion, ok bool, err error) {
	cookie, cookieErr := r.Cookie(sessionCookie)
	if cookieErr != nil {
		return service.UsuarioSesion{}, false, nil
	}
	user, err = h.auth.UsuarioDeSesion(r.Context(), cookie.Value)
	if errors.Is(err, service.ErrUnauthorized) {
		return service.UsuarioSesion{}, false, nil
	}
	if err != nil {
		return service.UsuarioSesion{}, false, err
	}
	return user, true, nil
}

// LoginPage muestra los dos accesos y, si viene ?rol=, el formulario de ese
// rol. Quien ya tiene sesión va directo a su inicio.
func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if user, ok, err := h.currentUser(r); err != nil {
		log.Printf("no se pudo verificar la sesión: %v", err)
		http.Error(w, "No se pudo verificar la sesión.", http.StatusInternalServerError)
		return
	} else if ok {
		http.Redirect(w, r, homeFor(user.Rol), http.StatusSeeOther)
		return
	}
	h.renderLogin(w, r, views.LoginData{Rol: validRol(r.URL.Query().Get("rol"))}, http.StatusOK)
}

// Login valida las credenciales para el rol elegido y abre la sesión.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario.", http.StatusBadRequest)
		return
	}
	rol := validRol(r.FormValue("rol"))
	data := views.LoginData{Rol: rol, Email: strings.TrimSpace(r.FormValue("email"))}
	if rol == "" {
		h.renderLogin(w, r, views.LoginData{Error: "Elige cómo quieres entrar."}, http.StatusBadRequest)
		return
	}

	session, err := h.auth.Login(r.Context(), data.Email, r.FormValue("contrasena"), rol)
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		data.Error = "Correo o contraseña incorrectos."
		h.renderLogin(w, r, data, http.StatusUnauthorized)
		return
	case errors.Is(err, service.ErrInvalid):
		data.Error = "Completa el correo y la contraseña."
		h.renderLogin(w, r, data, http.StatusBadRequest)
		return
	case err != nil:
		log.Printf("no se pudo iniciar sesión: %v", err)
		http.Error(w, "No se pudo iniciar sesión.", http.StatusInternalServerError)
		return
	}

	h.setSessionCookie(w, session.Token, session.Expira)
	http.Redirect(w, r, homeFor(session.Usuario.Rol), http.StatusSeeOther)
}

// Logout borra la sesión del servidor y de la cookie. Es POST para que un
// link ajeno no pueda cerrarle la sesión al usuario.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		if err := h.auth.Logout(r.Context(), cookie.Value); err != nil {
			log.Printf("no se pudo cerrar la sesión: %v", err)
		}
	}
	h.clearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func validRol(rol string) string {
	if rol == service.RolAlumno || rol == service.RolProfesor {
		return rol
	}
	return ""
}
