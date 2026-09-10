-- name: ListEjerciciosRutina :many
SELECT e.*
FROM rutina_ejercicio re
JOIN ejercicio e ON re.ejercicio_id = e.id
WHERE re.rutina_id = $1;

-- name: ClasesProfesor :many
SELECT
    u.nombre,
    u.apellido,
    c.nombre AS clase_nombre,
    c.max_alumnos,
    ch.dia_semana,
    ch.hora_inicio
FROM clase c
JOIN profesor p ON c.profesor_id = p.usuario_id
JOIN usuario u ON p.usuario_id = u.id
JOIN clase_horario ch ON ch.clase_id = c.id
WHERE c.profesor_id = $1;

-- name: ProfesoresdelDia :many
SELECT DISTINCT u.nombre, u.apellido
FROM profesor p
JOIN usuario u ON p.usuario_id = u.id
JOIN clase c ON c.profesor_id = p.usuario_id
JOIN clase_horario ch ON ch.clase_id = c.id
WHERE ch.dia_semana = $1;
