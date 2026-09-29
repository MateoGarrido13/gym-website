-- Rollback de clase_horario. inscripcion_clase debe bajarse antes.

ALTER TABLE clase_horario DROP CONSTRAINT IF EXISTS fk_clase_horario_clase;
DROP TABLE IF EXISTS clase_horario;
