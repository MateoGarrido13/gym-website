-- CONSULTAS DE RUTINA

-- name: CreateRutina :one
INSERT INTO rutina (nombre, duracion_semanas, profesor_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRutina :one
SELECT id, nombre, duracion_semanas, profesor_id
FROM rutina
WHERE id = $1;

-- name: ListRutinas :many
SELECT id, nombre, duracion_semanas, profesor_id
FROM rutina;

-- name: UpdateRutina :exec
UPDATE rutina
SET nombre = $1, duracion_semanas = $2
WHERE id = $3;

-- name: DeleteRutina :exec
DELETE FROM rutina
WHERE id = $1;

-- name: AddEjercicioARutina :exec
INSERT INTO rutina_ejercicio (rutina_id, ejercicio_id, orden, series, repeticiones, descanso_segundos)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DeleteEjercicioDeRutina :exec
DELETE FROM rutina_ejercicio
WHERE rutina_id = $1 AND ejercicio_id = $2;
