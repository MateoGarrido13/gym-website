CREATE TABLE inscripcion_clase (
    clase_horario_id  INT NOT NULL,
    alumno_id         INT NOT NULL,
    fecha_inscripcion  TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (clase_horario_id, alumno_id)
);

ALTER TABLE inscripcion_clase
    ADD CONSTRAINT fk_inscripcion_horario
        FOREIGN KEY (clase_horario_id) REFERENCES clase_horario(id) ON DELETE CASCADE;

ALTER TABLE inscripcion_clase
    ADD CONSTRAINT fk_inscripcion_alumno
        FOREIGN KEY (alumno_id) REFERENCES alumno(usuario_id) ON DELETE CASCADE;
