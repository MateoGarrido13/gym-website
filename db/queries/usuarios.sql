-- CONSULTAS DE USUARIOS 

-- name: CreateUsuario :one
INSERT INTO usuarios (email, contrasena, nombre, apellido, telefono, rol)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, nombre, apellido, telefono, rol, created_at;

-- name: GetUsuario :one
SELECT id, email, contrasena, nombre, apellido, telefono, rol, created_at
FROM usuarios
WHERE id = $1;

-- name: UpdateUsuario :exec
UPDATE usuarios
SET email = $1, contrasena = $2, nombre = $3, apellido = $4, telefono = $5
WHERE id = $6;

-- name: DeleteUsuario :exec
DELETE FROM usuarios
WHERE id = $1;
