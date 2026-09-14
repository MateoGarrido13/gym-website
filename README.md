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
