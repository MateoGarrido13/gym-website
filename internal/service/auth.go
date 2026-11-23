package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	db "PRACTICO_DOS/db/sqlc"
)

// RolAlumno, RolProfesor y RolAdmin son los valores de enum_rol.
const (
	RolAlumno   = "alumno"
	RolProfesor = "profesor"
	RolAdmin    = "administrador"
)

type AuthService struct {
	store AuthStore
}

func NewAuthService(store AuthStore) *AuthService {
	return &AuthService{store: store}
}

// hashFalso se compara cuando el email no existe, para que el tiempo de
// respuesta no revele si la cuenta existe.
var hashFalso, _ = hashPassword("contraseña-que-nadie-usa")

// Login valida email, contraseña y que el rol de la cuenta sea el elegido en
// la pantalla ("entrar como alumno" con una cuenta de profesor falla). Todos
// los fallos devuelven el mismo ErrUnauthorized para no filtrar cuál fue.
func (s *AuthService) Login(ctx context.Context, email, password, rol string) (SesionIniciada, error) {
	if rol != RolAlumno && rol != RolProfesor && rol != RolAdmin {
		return SesionIniciada{}, fmt.Errorf("%w: rol inválido", ErrInvalid)
	}
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return SesionIniciada{}, fmt.Errorf("%w: email y contraseña son obligatorios", ErrInvalid)
	}

	usuario, err := s.store.GetUsuarioByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		checkPassword(hashFalso, password)
		return SesionIniciada{}, ErrUnauthorized
	}
	if err != nil {
		return SesionIniciada{}, err
	}
	if !checkPassword(usuario.Contrasena, password) || usuario.Rol != rol {
		return SesionIniciada{}, ErrUnauthorized
	}

	token, err := nuevoToken()
	if err != nil {
		return SesionIniciada{}, err
	}
	expira := time.Now().Add(SesionDuracion)
	if err := s.store.CreateSesion(ctx, db.CreateSesionParams{
		ID:        hashToken(token),
		UsuarioID: usuario.ID,
		ExpiraEn:  expira,
	}); err != nil {
		return SesionIniciada{}, err
	}
	return SesionIniciada{
		Token:  token,
		Expira: expira,
		Usuario: UsuarioSesion{
			ID: usuario.ID, Email: usuario.Email, Nombre: usuario.Nombre,
			Apellido: usuario.Apellido, Rol: usuario.Rol,
		},
	}, nil
}

// UsuarioDeSesion resuelve el usuario dueño del token. Si el token no existe
// o la sesión expiró devuelve ErrUnauthorized.
func (s *AuthService) UsuarioDeSesion(ctx context.Context, token string) (UsuarioSesion, error) {
	if token == "" {
		return UsuarioSesion{}, ErrUnauthorized
	}
	row, err := s.store.GetSesionUsuario(ctx, hashToken(token))
	if errors.Is(err, sql.ErrNoRows) {
		return UsuarioSesion{}, ErrUnauthorized
	}
	if err != nil {
		return UsuarioSesion{}, err
	}
	return fromDBUsuarioSesion(row), nil
}

// Logout borra la sesión. Si el token no existe no es un error.
func (s *AuthService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	return s.store.DeleteSesion(ctx, hashToken(token))
}

// LimpiarExpiradas borra las sesiones vencidas.
func (s *AuthService) LimpiarExpiradas(ctx context.Context) error {
	return s.store.DeleteSesionesExpiradas(ctx)
}

// nuevoToken genera 32 bytes aleatorios (256 bits) en hexadecimal.
func nuevoToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
