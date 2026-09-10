CREATE TYPE enum_rol AS ENUM ('alumno', 'profesor');
CREATE TYPE enum_dia_semana AS ENUM (
    'lunes',
    'martes',
    'miercoles',
    'jueves',
    'viernes',
    'sabado',
    'domingo'
);

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

CREATE TABLE profesor (
    usuario_id    INT PRIMARY KEY,
    especialidad  VARCHAR(100),
    CONSTRAINT fk_profesor_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuario(id) ON DELETE CASCADE
);

CREATE TABLE ejercicio (
    id              SERIAL PRIMARY KEY,
    nombre          VARCHAR(100) NOT NULL,
    descripcion     TEXT,
    grupo_muscular  VARCHAR(100)
);

CREATE TABLE rutina (
    id                SERIAL PRIMARY KEY,
    nombre            VARCHAR(100) NOT NULL,
    duracion_semanas  INT NOT NULL,
    profesor_id       INT NOT NULL,
    CONSTRAINT fk_rutina_profesor
        FOREIGN KEY (profesor_id) REFERENCES profesor(usuario_id)
);

CREATE TABLE rutina_ejercicio (
    rutina_id          INT NOT NULL,
    ejercicio_id       INT NOT NULL,
    orden              INT NOT NULL,
    series             INT,
    repeticiones       INT,
    descanso_segundos  INT,
    PRIMARY KEY (rutina_id, ejercicio_id),
    CONSTRAINT fk_rutina_ejercicio_rutina
        FOREIGN KEY (rutina_id) REFERENCES rutina(id) ON DELETE CASCADE,
    CONSTRAINT fk_rutina_ejercicio_ejercicio
        FOREIGN KEY (ejercicio_id) REFERENCES ejercicio(id)
);

CREATE TABLE alumno (
    usuario_id         INT PRIMARY KEY,
    fecha_inscripcion  DATE NOT NULL,
    fecha_vto         DATE,
    tipo_plan         VARCHAR(50),
    rutina_id         INT,
    CONSTRAINT fk_alumno_usuario
        FOREIGN KEY (usuario_id) REFERENCES usuario(id) ON DELETE CASCADE,
    CONSTRAINT fk_alumno_rutina
        FOREIGN KEY (rutina_id) REFERENCES rutina(id)
);

CREATE TABLE clase (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL,
    max_alumnos INT NOT NULL,
    profesor_id INT NOT NULL,
    CONSTRAINT fk_clase_profesor
        FOREIGN KEY (profesor_id) REFERENCES profesor(usuario_id)
);

CREATE TABLE clase_horario (
    id          SERIAL PRIMARY KEY,
    clase_id    INT NOT NULL,
    dia_semana  enum_dia_semana NOT NULL,
    hora_inicio TIME NOT NULL,
    CONSTRAINT fk_clase_horario_clase
        FOREIGN KEY (clase_id) REFERENCES clase(id) ON DELETE CASCADE
);

CREATE TABLE inscripcion_clase (
    clase_horario_id  INT NOT NULL,
    alumno_id         INT NOT NULL,
    fecha_inscripcion  TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (clase_horario_id, alumno_id),
    CONSTRAINT fk_inscripcion_horario
        FOREIGN KEY (clase_horario_id) REFERENCES clase_horario(id) ON DELETE CASCADE,
    CONSTRAINT fk_inscripcion_alumno
        FOREIGN KEY (alumno_id) REFERENCES alumno(usuario_id) ON DELETE CASCADE
);
