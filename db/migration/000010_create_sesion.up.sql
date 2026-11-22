-- Sesiones de login. "id" guarda el SHA-256 (hex) del token que viaja en la
-- cookie, nunca el token en sí: si se filtra la tabla, no sirve para entrar.
CREATE TABLE sesion (
    id          CHAR(64) PRIMARY KEY,
    usuario_id  INT NOT NULL,
    expira_en   TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Al borrar un usuario se cierran todas sus sesiones.
ALTER TABLE sesion
    ADD CONSTRAINT fk_sesion_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuario(id) ON DELETE CASCADE;

CREATE INDEX idx_sesion_usuario ON sesion(usuario_id);
