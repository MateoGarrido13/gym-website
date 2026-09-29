CREATE TABLE alumno (
    usuario_id         INT PRIMARY KEY,
    fecha_inscripcion  DATE NOT NULL,
    fecha_vto          DATE,
    tipo_plan          VARCHAR(50),
    rutina_id          INT
);

ALTER TABLE alumno
    ADD CONSTRAINT fk_alumno_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuario(id) ON DELETE CASCADE;

ALTER TABLE alumno
    ADD CONSTRAINT fk_alumno_rutina
        FOREIGN KEY (rutina_id) REFERENCES rutina(id);
