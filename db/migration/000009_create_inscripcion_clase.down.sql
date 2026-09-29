ALTER TABLE inscripcion_clase DROP CONSTRAINT IF EXISTS fk_inscripcion_alumno;
ALTER TABLE inscripcion_clase DROP CONSTRAINT IF EXISTS fk_inscripcion_horario;
DROP TABLE IF EXISTS inscripcion_clase;
