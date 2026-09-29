-- Rollback de clase. clase_horario debe bajarse antes.

ALTER TABLE clase DROP CONSTRAINT IF EXISTS fk_clase_profesor;
DROP TABLE IF EXISTS clase;
