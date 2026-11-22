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

-- Rutina asignada a un alumno, con el nombre de su profesor. Sin filas si el
-- alumno no tiene rutina.
-- name: GetRutinaDeAlumno :one
SELECT r.id, r.nombre, r.duracion_semanas,
       u.nombre AS profesor_nombre, u.apellido AS profesor_apellido
FROM alumno a
JOIN rutina r ON r.id = a.rutina_id
JOIN usuario u ON u.id = r.profesor_id
WHERE a.usuario_id = $1;

-- name: ListEjerciciosDetalleDeRutina :many
SELECT re.orden, re.series, re.repeticiones, re.descanso_segundos,
       e.nombre, e.descripcion, e.grupo_muscular
FROM rutina_ejercicio re
JOIN ejercicio e ON e.id = re.ejercicio_id
WHERE re.rutina_id = $1
ORDER BY re.orden;
