-- CONSULTAS DE ALUMNO

-- name: InsertAlumno :exec
INSERT INTO alumno (usuario_id, fecha_inscripcion, fecha_vto, tipo_plan, rutina_id)
VALUES ($1, $2, $3, $4, $5);

-- name: GetAlumnoCompleto :one
SELECT
    u.id, u.email, u.nombre, u.apellido, u.telefono, u.rol, u.created_at,
    a.fecha_inscripcion, a.fecha_vto, a.tipo_plan, a.rutina_id
FROM usuario u
JOIN alumno a ON u.id = a.usuario_id
WHERE u.id = $1;

-- name: ListAlumnosCompletos :many
SELECT
    u.id, u.email, u.nombre, u.apellido, u.telefono, u.rol, u.created_at,
    a.fecha_inscripcion, a.fecha_vto, a.tipo_plan, a.rutina_id
FROM usuario u
JOIN alumno a ON u.id = a.usuario_id;

-- name: UpdateAlumnoDetalles :exec
UPDATE alumno
SET fecha_vto = $1, tipo_plan = $2, rutina_id = $3
WHERE usuario_id = $4;

-- el alumno sera eliminado desde la tabla 'usuario'
