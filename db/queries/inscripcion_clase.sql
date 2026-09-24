-- name: ListInscriptos :many
SELECT ic.alumno_id, u.nombre, u.apellido
FROM inscripcion_clase ic
JOIN alumno a ON ic.alumno_id = a.usuario_id
JOIN usuario u ON a.usuario_id = u.id
WHERE ic.clase_horario_id = $1;

-- name: CreateInscripcion :one
INSERT INTO inscripcion_clase (clase_horario_id, alumno_id)
VALUES ($1, $2)
RETURNING clase_horario_id, alumno_id, fecha_inscripcion;

-- name: DeleteInscripcion :exec
DELETE FROM inscripcion_clase
WHERE clase_horario_id = $1 AND alumno_id = $2;

-- name: GetCupoDeHorario :one
SELECT c.max_alumnos,
       COUNT(ic.alumno_id)::int AS cantidad_inscriptos
FROM clase_horario ch
JOIN clase c ON c.id = ch.clase_id
LEFT JOIN inscripcion_clase ic ON ic.clase_horario_id = ch.id
WHERE ch.id = $1
GROUP BY c.max_alumnos;
