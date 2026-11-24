-- Revertir el enum de tipo_plan. PostgreSQL no permite quitar valores de un
-- enum; para revertir completamente hay que borrar los alumnos y recrear el
-- tipo, o hacer `make reset`.

UPDATE alumno SET tipo_plan = NULL;

ALTER TABLE alumno
    ALTER COLUMN tipo_plan TYPE VARCHAR(50);

DROP TYPE IF EXISTS enum_tipo_plan;
