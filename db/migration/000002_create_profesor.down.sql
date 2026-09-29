-- Rollback de profesor. rutina y clase deben bajarse antes.

ALTER TABLE profesor DROP CONSTRAINT IF EXISTS fk_profesor_usuario;
DROP TABLE IF EXISTS profesor;
