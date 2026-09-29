CREATE TABLE usuario (
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(255) NOT NULL UNIQUE,
    contrasena  VARCHAR(255) NOT NULL,
    nombre      VARCHAR(100) NOT NULL,
    apellido    VARCHAR(100) NOT NULL,
    telefono    VARCHAR(30),
    rol         enum_rol NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
