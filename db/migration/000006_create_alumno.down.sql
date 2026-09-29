-- Rollback de alumno. inscripcion_clase debe bajarse antes.

ALTER TABLE alumno DROP CONSTRAINT IF EXISTS fk_alumno_rutina;
ALTER TABLE alumno DROP CONSTRAINT IF EXISTS fk_alumno_usuario;
DROP TABLE IF EXISTS alumno;
