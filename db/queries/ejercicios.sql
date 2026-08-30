-- CONSULTAS DE EJERCICIOS 

-- name: CreateEjercicio :one
INSERT INTO ejercicios (nombre, descripcion)
VALUES ($1, $2)
RETURNING id, nombre, descripcion;

-- name: GetEjercicio :one
SELECT id, nombre, descripcion
FROM ejercicios
WHERE id = $1;

-- name: ListEjercicios :many
SELECT id, nombre, descripcion
FROM ejercicios;

-- name: UpdateEjercicio :exec
UPDATE ejercicios
SET nombre = $1, descripcion = $2
WHERE id = $3;

-- name: DeleteEjercicio :exec
DELETE FROM ejercicios
WHERE id = $1;
