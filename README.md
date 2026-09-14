# gym-website

Página web para gimnasios. El objetivo es digitalizar las rutinas, agilizar su actualización y canalizar la inscripción a clases e información en la aplicación. A su vez, mejorar la cercanía con el cliente y su experiencia.

Este repositorio corresponde al TP2: persistencia en PostgreSQL, consultas generadas con sqlc y tests de integración.

## Cómo ejecutar

Requisitos: Go, Docker y Docker Compose.

```bash
git clone <url-del-repo>
cd gym-website
git checkout tp2
make test
```

`make test` hace, en orden:

1. Instala sqlc (si hace falta) y genera el código en `db/sqlc/`
2. Compila `main.go`
3. Baja contenedores y volúmenes previos
4. Levanta Postgres y espera a que esté listo (healthcheck)
5. Corre los tests (`go test ./test/...`)
6. Baja contenedores y volúmenes

## Persistencia

Los datos viven en **PostgreSQL 15** (`postgres:15-alpine`), definido en `docker-compose.yml`.

- Base: `tp2_db`
- Usuario / contraseña: `postgres` / `postgres`
- Puerto: `5432`
- Volumen: `postgres_data`

El esquema está en `db/schema/schema.sql`. Compose lo monta en `/docker-entrypoint-initdb.d/`, así que Postgres lo aplica **solo la primera vez** que se crea el volumen. Por eso el Makefile usa `docker compose down -v` antes y después: cada corrida arranca con el schema limpio.

Las consultas SQL están en `db/queries/`. `sqlc.yaml` las traduce a Go en `db/sqlc/` (paquete `db`). Ese código no se edita a mano.

La app y los tests se conectan con el driver `github.com/lib/pq`. Los tests usan:

`host=localhost port=5432 user=postgres password=postgres dbname=tp2_db sslmode=disable`

No hay migraciones incrementales: un solo schema al init del contenedor.

### Modelo

`usuario` es la identidad (email, contraseña, nombre, apellido, teléfono, `rol`). `alumno` y `profesor` la especializan 1:1 (`usuario_id` es PK y FK). Borrar un usuario, por `ON DELETE CASCADE`, borra el alumno o profesor asociado.

| Tabla | Qué guarda |
|---|---|
| `usuario` | Datos personales y rol (`alumno` / `profesor`) |
| `alumno` | Fecha de inscripción, vencimiento, tipo de plan, rutina asignada |
| `profesor` | Especialidad |
| `ejercicio` | Nombre, descripción, grupo muscular |
| `rutina` | Nombre, duración en semanas, profesor |
| `rutina_ejercicio` | Ejercicios de una rutina (orden, series, reps, descanso) |
| `clase` | Nombre, cupo, profesor |
| `clase_horario` | Día de la semana y hora |
| `inscripcion_clase` | Alumno inscripto a un horario |

Consultas cubiertas: CRUD de usuario, ejercicio, profesor, rutina y clase; alta/listado/update de alumno (join con usuario); inscripciones; joins (ejercicios de una rutina, clases de un profesor, profesores de un día).

## Tests

Paquete `testing` en `test/`. Prueban la persistencia contra la base real: usuario, alumno, profesor, ejercicio, rutina, clase, inscripciones y consultas multitabla. Cada `make test` deja la DB en un estado fresco.

## Funcionalidades previstas del sistema

**Usuario profesor**

- Carga rutinas con ejercicios
- Asigna las rutinas a los alumnos
- Gestionar clases (pack de clases)
- Dashboard con métricas por persona (pesos, concurrencia, ritmos)

**Usuario alumno**

- Ve la rutina que le fue asignada
- Accede a los ejercicios (imágenes/videos)
- Ver calendario e inscribirse a una clase

En este TP la capa de persistencia y los tests cubren el modelo de datos; `main.go` sirve archivos estáticos y todavía no usa la base.

## Requisitos / versiones

- Go `1.26.5`
- sqlc `v1.31.1`
- PostgreSQL `15-alpine` (imagen Docker)
- Docker Compose `3.8`
- `github.com/lib/pq` `v1.12.3`

## ADICIONAL

Además de la persistencia y los tests del TP, **en la rama tp3 hay una API HTTP de muestra** para ejercicios y alumnos. Quisimos aprozimarnos a una implementacion mas prolija con handlers y servicios, en lugar de llamarlo solo desde el paquete `testing`.No reemplaza `make test`: sirve para ver el mismo `*db.Queries` usado desde un servidor.

La utilidad es tener un  patrón de capas (handler → servicio → puerto; `*db.Queries` implementa el puerto) y comprobar el modelo por HTTP. Ejercicios cubren un CRUD directo. Alumnos cubren el alta compuesta usuario+alumno (`CreateAlumno`) y el borrado por cascade sobre `usuario`. El resto de tablas sigue cubierto por sqlc y por los tests de integración.

Cómo probarla (Postgres tiene que estar arriba; `make test` la baja al terminar):

```bash
make up
make run
```

El servidor queda en `http://localhost:8080`. `GET /` sirve `static/`. Las rutas JSON están bajo `/api/ejercicios` y `/api/alumnos` (ver la tabla en **API (muestra)**).

```bash
curl -s http://localhost:8080/
curl -s http://localhost:8080/api/ejercicios
curl -s -X POST http://localhost:8080/api/ejercicios \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"sentadilla","descripcion":"barra","grupo_muscular":"piernas"}'
curl -s -X POST http://localhost:8080/api/alumnos \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@example.com","contrasena":"secret123","nombre":"Ana","apellido":"Lopez","fecha_inscripcion":"2026-09-10","tipo_plan":"mensual"}'
```

Si `make up` falla porque el puerto `5432` ya está ocupado, alcanza con que esa instancia tenga la base `tp2_db` y el schema; `make run` usa la misma cadena que los tests. Al terminar: `make down` (solo si levantaste Compose en este paso).
