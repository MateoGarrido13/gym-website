-- Datos de desarrollo para poder entrar a la app apenas se levanta.
-- NO es una migración: no lo aplica init.sh ni los tests. Lo carga
-- `make seed` (y `make run-app`). Es idempotente: si ya existe el profesor
-- de prueba, no hace nada.
--
-- Usuarios (contraseñas solo para desarrollo, guardadas con bcrypt):
--   profe@gym.test      profe1234    profesor
--   alumno@gym.test     alumno1234   alumno con rutina asignada
--   sinrutina@gym.test  alumno1234   alumno sin rutina

DO $$
DECLARE
    v_profe   INT;
    v_rutina  INT;
    v_alumno  INT;
    v_ej      INT;
BEGIN
    IF EXISTS (SELECT 1 FROM usuario WHERE email = 'profe@gym.test') THEN
        RAISE NOTICE 'seed: los datos de prueba ya existen, no se hace nada';
        RETURN;
    END IF;

    INSERT INTO usuario (email, contrasena, nombre, apellido, rol)
    VALUES ('profe@gym.test', '$2a$10$wDIASYgaxVtxwupwQ5oYpusSx2iIClegYB8PBUPGPx0ArclGimxAu', 'Pablo', 'Profesor', 'profesor')
    RETURNING id INTO v_profe;
    INSERT INTO profesor (usuario_id, especialidad) VALUES (v_profe, 'Fuerza');

    INSERT INTO rutina (nombre, duracion_semanas, profesor_id)
    VALUES ('Fuerza full body', 8, v_profe)
    RETURNING id INTO v_rutina;

    INSERT INTO ejercicio (nombre, descripcion, grupo_muscular)
    VALUES ('Sentadilla', 'Barra apoyada en la espalda alta, bajar hasta paralelo.', 'Piernas')
    RETURNING id INTO v_ej;
    INSERT INTO rutina_ejercicio (rutina_id, ejercicio_id, orden, series, repeticiones, descanso_segundos)
    VALUES (v_rutina, v_ej, 1, 4, 12, 60);

    INSERT INTO ejercicio (nombre, descripcion, grupo_muscular)
    VALUES ('Press de banca', 'Acostado en el banco, bajar la barra al pecho y empujar.', 'Pecho')
    RETURNING id INTO v_ej;
    INSERT INTO rutina_ejercicio (rutina_id, ejercicio_id, orden, series, repeticiones, descanso_segundos)
    VALUES (v_rutina, v_ej, 2, 4, 10, 75);

    INSERT INTO ejercicio (nombre, descripcion, grupo_muscular)
    VALUES ('Remo con barra', 'Torso inclinado, llevar la barra hacia el abdomen.', 'Espalda')
    RETURNING id INTO v_ej;
    INSERT INTO rutina_ejercicio (rutina_id, ejercicio_id, orden, series, repeticiones, descanso_segundos)
    VALUES (v_rutina, v_ej, 3, 4, 8, 90);

    INSERT INTO ejercicio (nombre, descripcion, grupo_muscular)
    VALUES ('Plancha', 'Apoyo en antebrazos, cuerpo recto. Se mide en segundos.', 'Abdomen')
    RETURNING id INTO v_ej;
    INSERT INTO rutina_ejercicio (rutina_id, ejercicio_id, orden, series, descanso_segundos)
    VALUES (v_rutina, v_ej, 4, 3, 45);

    INSERT INTO usuario (email, contrasena, nombre, apellido, rol)
    VALUES ('alumno@gym.test', '$2a$10$1/zDrLMJ6Ijw75rcF2XhVe3hZvXAE.PA7QWLSWG1u/w4Abl3H0OL2', 'Ana', 'Alumna', 'alumno')
    RETURNING id INTO v_alumno;
    INSERT INTO alumno (usuario_id, fecha_inscripcion, tipo_plan, rutina_id)
    VALUES (v_alumno, CURRENT_DATE, 'Mensual', v_rutina);

    INSERT INTO usuario (email, contrasena, nombre, apellido, rol)
    VALUES ('sinrutina@gym.test', '$2a$10$1/zDrLMJ6Ijw75rcF2XhVe3hZvXAE.PA7QWLSWG1u/w4Abl3H0OL2', 'Beto', 'Sinrutina', 'alumno')
    RETURNING id INTO v_alumno;
    INSERT INTO alumno (usuario_id, fecha_inscripcion, tipo_plan)
    VALUES (v_alumno, CURRENT_DATE, 'Mensual');
END $$;
