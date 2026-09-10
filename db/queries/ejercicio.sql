-- CONSULTAS DE EJERCICIO

-- name: CreateEjercicio :one
INSERT INTO ejercicio (nombre, descripcion, grupo_muscular)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEjercicio :one
SELECT id, nombre, descripcion, grupo_muscular
FROM ejercicio
WHERE id = $1;

-- name: ListEjercicios :many
SELECT id, nombre, descripcion, grupo_muscular
FROM ejercicio;

-- name: UpdateEjercicio :exec
UPDATE ejercicio
SET nombre = $1, descripcion = $2, grupo_muscular = $3
WHERE id = $4;

-- name: DeleteEjercicio :exec
DELETE FROM ejercicio
WHERE id = $1;
