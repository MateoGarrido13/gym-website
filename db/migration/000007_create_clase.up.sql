CREATE TABLE clase (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL,
    max_alumnos INT NOT NULL,
    profesor_id INT NOT NULL
);

ALTER TABLE clase
    ADD CONSTRAINT fk_clase_profesor
        FOREIGN KEY (profesor_id) REFERENCES profesor(usuario_id);
