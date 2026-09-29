ALTER TABLE rutina_ejercicio DROP CONSTRAINT IF EXISTS fk_rutina_ejercicio_ejercicio;
ALTER TABLE rutina_ejercicio DROP CONSTRAINT IF EXISTS fk_rutina_ejercicio_rutina;
DROP TABLE IF EXISTS rutina_ejercicio;
