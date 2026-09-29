CREATE TABLE ejercicio (
    id              SERIAL PRIMARY KEY,
    nombre          VARCHAR(100) NOT NULL,
    descripcion     TEXT,
    grupo_muscular  VARCHAR(100)
);
