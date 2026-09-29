CREATE TABLE profesor (
    usuario_id    INT PRIMARY KEY,
    especialidad  VARCHAR(100)
);

ALTER TABLE profesor
    ADD CONSTRAINT fk_profesor_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuario(id) ON DELETE CASCADE;
