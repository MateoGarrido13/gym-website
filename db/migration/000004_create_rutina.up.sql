CREATE TABLE rutina (
    id                SERIAL PRIMARY KEY,
    nombre            VARCHAR(100) NOT NULL,
    duracion_semanas  INT NOT NULL,
    profesor_id       INT NOT NULL
);

ALTER TABLE rutina
    ADD CONSTRAINT fk_rutina_profesor
        FOREIGN KEY (profesor_id) REFERENCES profesor(usuario_id);
