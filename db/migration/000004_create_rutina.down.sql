-- Rollback de rutina. rutina_ejercicio y alumno deben bajarse antes.

ALTER TABLE rutina DROP CONSTRAINT IF EXISTS fk_rutina_profesor;
DROP TABLE IF EXISTS rutina;
