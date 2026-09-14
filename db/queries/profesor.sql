-- CONSULTAS DE PROFESOR

-- name: InsertProfesor :exec
INSERT INTO profesor (usuario_id, especialidad)
VALUES ($1, $2);

-- name: GetProfesorCompleto :one
SELECT
    u.id, u.email, u.nombre, u.apellido, u.telefono, u.rol, u.created_at,
    p.especialidad
FROM usuario u
JOIN profesor p ON u.id = p.usuario_id
WHERE u.id = $1;

-- name: ListProfesoresCompletos :many
SELECT
    u.id, u.email, u.nombre, u.apellido, u.telefono, u.rol, u.created_at,
    p.especialidad
FROM usuario u
JOIN profesor p ON u.id = p.usuario_id;

-- name: UpdateProfesor :exec
UPDATE profesor
SET especialidad = $1
WHERE usuario_id = $2;

-- el profesor sera eliminado desde la tabla 'usuario'
