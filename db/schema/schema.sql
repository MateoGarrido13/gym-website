-- USUARIOS (Datos de acceso y personales comunes)
CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    contrasena VARCHAR(255) NOT NULL, -- Guardará el hash de la clave por seguridad
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    telefono VARCHAR(50),
    rol VARCHAR(20) NOT NULL CHECK (rol IN ('profesor', 'alumno')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- EJERCICIOS (Catálogo general de ejercicios)
CREATE TABLE ejercicios (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    descripcion TEXT
);

-- RUTINAS (Planes de entrenamiento)
CREATE TABLE rutinas (
    id SERIAL PRIMARY KEY,
    duracion_semanas INT NOT NULL DEFAULT 4
);

-- TABLA INTERMEDIA: RUTINA_EJERCICIOS (Relación muchos a muchos)
CREATE TABLE rutina_ejercicios (
    rutina_id INT REFERENCES rutinas(id) ON DELETE CASCADE,
    ejercicio_id INT REFERENCES ejercicios(id) ON DELETE CASCADE,
    PRIMARY KEY (rutina_id, ejercicio_id)
);

-- ALUMNOS (Extiende la tabla usuarios)
CREATE TABLE alumnos (
    usuario_id INT PRIMARY KEY REFERENCES usuarios(id) ON DELETE CASCADE,
    fecha_inscripcion DATE NOT NULL DEFAULT CURRENT_DATE,
    fecha_vto DATE,
    tipo_plan VARCHAR(100),
    rutina_id INT REFERENCES rutinas(id) ON DELETE SET NULL
);

-- PROFESORES (Extiende la tabla usuarios)
CREATE TABLE profesores (
    usuario_id INT PRIMARY KEY REFERENCES usuarios(id) ON DELETE CASCADE,
    horario VARCHAR(255), -- Horarios de disponibilidad (ej. "Turno Mañana")
    especialidad VARCHAR(255)
);

-- CLASES (Agenda de actividades grupales)
CREATE TABLE clases (
    id SERIAL PRIMARY KEY,
    dias VARCHAR(20) NOT NULL, -- Ej: "Lunes,Miércoles,Viernes"
    hora TIME NOT NULL,       -- Hora de inicio
    max_alumnos INT NOT NULL,
    profesor_id INT REFERENCES profesores(usuario_id) ON DELETE SET NULL
);

-- TABLA INTERMEDIA: INSCRIPCIONES_CLASES (Alumnos anotados en clases)
CREATE TABLE inscripciones_clases (
    clase_id INT REFERENCES clases(id) ON DELETE CASCADE,
    alumno_id INT REFERENCES alumnos(usuario_id) ON DELETE CASCADE,
    PRIMARY KEY (clase_id, alumno_id),
    dia VARCHAR(20) PRIMARY KEY , -- Seria conveniente tagearlas como LUN,MAR,MIE,JUE,VIE
    hora TIME NOT NULL
);
-- Posible inconsistencia entre que el dia registrado para esta clase no sea el permitido por la tabla CLASE
