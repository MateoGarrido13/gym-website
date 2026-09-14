-- CONSULTAS DE CLASE

-- name: CreateClase :one
INSERT INTO clase (nombre, max_alumnos, profesor_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetClase :one
SELECT id, nombre, max_alumnos, profesor_id
FROM clase
WHERE id = $1;

-- name: ListClases :many
SELECT id, nombre, max_alumnos, profesor_id
FROM clase;

-- name: UpdateClase :exec
UPDATE clase
SET nombre = $1, max_alumnos = $2
WHERE id = $3;

-- name: DeleteClase :exec
DELETE FROM clase
WHERE id = $1;

-- name: CreateClaseHorario :one
INSERT INTO clase_horario (clase_id, dia_semana, hora_inicio)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetClaseHorario :one
SELECT id, clase_id, dia_semana, hora_inicio
FROM clase_horario
WHERE id = $1;

-- name: ListHorariosDeClase :many
SELECT id, clase_id, dia_semana, hora_inicio
FROM clase_horario
WHERE clase_id = $1;

-- name: DeleteClaseHorario :exec
DELETE FROM clase_horario
WHERE id = $1;
