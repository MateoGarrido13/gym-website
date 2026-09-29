-- Tabla puente rutina ↔ ejercicio. Sale de schema.sql junto con esas entidades.

CREATE TABLE rutina_ejercicio (
    rutina_id          INT NOT NULL,
    ejercicio_id       INT NOT NULL,
    orden              INT NOT NULL,
    series             INT,
    repeticiones       INT,
    descanso_segundos  INT,
    PRIMARY KEY (rutina_id, ejercicio_id)
);

ALTER TABLE rutina_ejercicio
    ADD CONSTRAINT fk_rutina_ejercicio_rutina
        FOREIGN KEY (rutina_id) REFERENCES rutina(id) ON DELETE CASCADE;

ALTER TABLE rutina_ejercicio
    ADD CONSTRAINT fk_rutina_ejercicio_ejercicio
        FOREIGN KEY (ejercicio_id) REFERENCES ejercicio(id);
