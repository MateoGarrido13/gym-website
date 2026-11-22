-- CONSULTAS DE SESION

-- name: CreateSesion :exec
INSERT INTO sesion (id, usuario_id, expira_en)
VALUES ($1, $2, $3);

-- Devuelve el usuario dueño de la sesión solo si todavía no expiró.
-- name: GetSesionUsuario :one
SELECT u.id, u.email, u.nombre, u.apellido, u.rol, s.expira_en
FROM sesion s
JOIN usuario u ON u.id = s.usuario_id
WHERE s.id = $1 AND s.expira_en > NOW();

-- name: DeleteSesion :exec
DELETE FROM sesion
WHERE id = $1;

-- name: DeleteSesionesExpiradas :exec
DELETE FROM sesion
WHERE expira_en <= NOW();
