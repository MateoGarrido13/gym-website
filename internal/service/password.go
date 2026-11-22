package service

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword devuelve el hash bcrypt de la contraseña. Las contraseñas no se
// recortan: los espacios cuentan. bcrypt rechaza las de más de 72 bytes.
func hashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return "", fmt.Errorf("%w: la contraseña no puede superar los 72 bytes", ErrInvalid)
	}
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
