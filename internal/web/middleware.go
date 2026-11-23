package web

import (
	"context"
	"log"
	"net/http"
	"slices"

	"PRACTICO_DOS/internal/service"
)

type ctxKey struct{}

// UserFrom devuelve el usuario que RequireRole dejó en el contexto.
func UserFrom(ctx context.Context) (service.UsuarioSesion, bool) {
	user, ok := ctx.Value(ctxKey{}).(service.UsuarioSesion)
	return user, ok
}

// RequireRole deja pasar solo a quien tenga sesión válida con alguno de esos
// roles y le pasa el usuario al handler por el contexto.
//   - sin sesión → /login
//   - sesión de otro rol → al inicio de su propio rol
//
// Si el pedido viene de HTMX, un redirect normal metería la página de login
// dentro del fragmento; por eso se responde con HX-Redirect.
func (h *Handler) RequireRole(next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok, err := h.currentUser(r)
		if err != nil {
			log.Printf("no se pudo verificar la sesión: %v", err)
			http.Error(w, "No se pudo verificar la sesión.", http.StatusInternalServerError)
			return
		}
		if !ok {
			h.clearSessionCookie(w) // por si la cookie quedó con un token vencido
			redirect(w, r, "/login")
			return
		}
		if !slices.Contains(roles, user.Rol) {
			redirect(w, r, homeFor(user.Rol))
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, user)))
	}
}

func redirect(w http.ResponseWriter, r *http.Request, to string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
}
