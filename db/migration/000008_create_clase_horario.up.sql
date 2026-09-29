-- Horarios de una clase. Necesaria antes de inscripcion_clase.

CREATE TABLE clase_horario (
    id          SERIAL PRIMARY KEY,
    clase_id    INT NOT NULL,
    dia_semana  enum_dia_semana NOT NULL,
    hora_inicio TIME NOT NULL
);

ALTER TABLE clase_horario
    ADD CONSTRAINT fk_clase_horario_clase
        FOREIGN KEY (clase_id) REFERENCES clase(id) ON DELETE CASCADE;
