-- Tipo de plan de alumno: Mensual, Trimestral o Anual.

CREATE TYPE enum_tipo_plan AS ENUM (
    'Mensual',
    'Trimestral',
    'Anual'
);

ALTER TABLE alumno
    ALTER COLUMN tipo_plan TYPE enum_tipo_plan
    USING tipo_plan::enum_tipo_plan;
