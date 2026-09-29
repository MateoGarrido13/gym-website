-- Tipos de PostgreSQL. Van en la migración 000000 para existir antes de
-- usuario (enum_rol) y clase_horario (enum_dia_semana).

CREATE TYPE enum_rol AS ENUM ('alumno', 'profesor');
CREATE TYPE enum_dia_semana AS ENUM (
    'lunes',
    'martes',
    'miercoles',
    'jueves',
    'viernes',
    'sabado',
    'domingo'
);
