# =============================================================================
#    Healthcheck + variable por defecto
#    Define `base_url` (disponible para todas las requests siguientes) y
#    verifica que el servidor responde.
# =============================================================================
GET {{base_url}}/api/ejercicios
[Options]
variable: base_url=http://localhost:8080
HTTP 200


# =============================================================================
# 1. PROFESOR  (alta compuesta: usuario + profesor)
# =============================================================================

# 1.1 Crear profesor -> 201, captura del id autogenerado
POST {{base_url}}/api/profesores
Content-Type: application/json
{
  "email": "profe@example.com",
  "contrasena": "secret123",
  "nombre": "Juan",
  "apellido": "Perez",
  "telefono": "1130000000",
  "especialidad": "musculacion"
}
HTTP 201
[Captures]
profesor_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.email" == "profe@example.com"
jsonpath "$.nombre" == "Juan"
jsonpath "$.rol" == "profesor"
jsonpath "$.especialidad" == "musculacion"

# 1.2 Listar profesores -> 200 y arreglo no vacío
GET {{base_url}}/api/profesores
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1
jsonpath "$[0].id" isInteger

# 1.3 Obtener el profesor por id dinámico -> 200 y campos coincidentes
GET {{base_url}}/api/profesores/{{profesor_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{profesor_id}}
jsonpath "$.email" == "profe@example.com"
jsonpath "$.nombre" == "Juan"
jsonpath "$.apellido" == "Perez"

# 1.4 Actualizar (solo cambia la especialidad) -> 200 con el cambio aplicado
PUT {{base_url}}/api/profesores/{{profesor_id}}
Content-Type: application/json
{
  "especialidad": "crossfit"
}
HTTP 200
[Asserts]
jsonpath "$.id" == {{profesor_id}}
jsonpath "$.especialidad" == "crossfit"


# =============================================================================
# 2. EJERCICIO  (CRUD directo)
# =============================================================================

# 2.1 Crear ejercicio -> 201, captura del id
POST {{base_url}}/api/ejercicios
Content-Type: application/json
{
  "nombre": "sentadilla",
  "descripcion": "barra",
  "grupo_muscular": "piernas"
}
HTTP 201
[Captures]
ejercicio_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.nombre" == "sentadilla"
jsonpath "$.grupo_muscular" == "piernas"

# 2.2 Listar ejercicios -> 200 y arreglo
GET {{base_url}}/api/ejercicios
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1

# 2.3 Obtener por id dinámico
GET {{base_url}}/api/ejercicios/{{ejercicio_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{ejercicio_id}}
jsonpath "$.nombre" == "sentadilla"

# 2.4 Actualizar -> 200 con el nombre nuevo
PUT {{base_url}}/api/ejercicios/{{ejercicio_id}}
Content-Type: application/json
{
  "nombre": "sentadilla frontal",
  "descripcion": "barra",
  "grupo_muscular": "piernas"
}
HTTP 200
[Asserts]
jsonpath "$.nombre" == "sentadilla frontal"
jsonpath "$.grupo_muscular" == "piernas"


# =============================================================================
# 3. RUTINA  (depende del profesor)
# =============================================================================

# 3.1 Crear rutina -> 201, captura del id
POST {{base_url}}/api/rutinas
Content-Type: application/json
{
  "nombre": "Full body",
  "duracion_semanas": 8,
  "profesor_id": {{profesor_id}}
}
HTTP 201
[Captures]
rutina_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.nombre" == "Full body"
jsonpath "$.duracion_semanas" == 8
jsonpath "$.profesor_id" == {{profesor_id}}

# 3.2 Listar rutinas
GET {{base_url}}/api/rutinas
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1

# 3.3 Obtener por id dinámico
GET {{base_url}}/api/rutinas/{{rutina_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{rutina_id}}
jsonpath "$.profesor_id" == {{profesor_id}}

# 3.4 Actualizar -> 200
PUT {{base_url}}/api/rutinas/{{rutina_id}}
Content-Type: application/json
{
  "nombre": "Full body v2",
  "duracion_semanas": 10
}
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Full body v2"
jsonpath "$.duracion_semanas" == 10


# =============================================================================
# 4. RUTINA_EJERCICIO  (relación anidada rutina <-> ejercicio)
# =============================================================================

# 4.1 Agregar ejercicio a la rutina -> 201
POST {{base_url}}/api/rutinas/{{rutina_id}}/ejercicios
Content-Type: application/json
{
  "ejercicio_id": {{ejercicio_id}},
  "orden": 1,
  "series": 4,
  "repeticiones": 10,
  "descanso_segundos": 90
}
HTTP 201
[Asserts]
jsonpath "$.rutina_id" == {{rutina_id}}
jsonpath "$.ejercicio_id" == {{ejercicio_id}}
jsonpath "$.orden" == 1
jsonpath "$.series" == 4
jsonpath "$.repeticiones" == 10

# 4.2 Listar ejercicios de la rutina -> 200 y no vacío
GET {{base_url}}/api/rutinas/{{rutina_id}}/ejercicios
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1
jsonpath "$[0].ejercicio_id" == {{ejercicio_id}}

# 4.3 Actualizar orden/series/reps -> 200
PUT {{base_url}}/api/rutinas/{{rutina_id}}/ejercicios/{{ejercicio_id}}
Content-Type: application/json
{
  "orden": 1,
  "series": 5,
  "repeticiones": 8,
  "descanso_segundos": 120
}
HTTP 200
[Asserts]
jsonpath "$.series" == 5
jsonpath "$.repeticiones" == 8
jsonpath "$.descanso_segundos" == 120


# =============================================================================
# 5. CLASE  (depende del profesor)
# =============================================================================

# 5.1 Crear clase -> 201, captura del id
POST {{base_url}}/api/clases
Content-Type: application/json
{
  "nombre": "Funcional",
  "max_alumnos": 20,
  "profesor_id": {{profesor_id}}
}
HTTP 201
[Captures]
clase_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.nombre" == "Funcional"
jsonpath "$.max_alumnos" == 20
jsonpath "$.profesor_id" == {{profesor_id}}

# 5.2 Listar clases
GET {{base_url}}/api/clases
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1

# 5.3 Obtener por id dinámico
GET {{base_url}}/api/clases/{{clase_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{clase_id}}
jsonpath "$.nombre" == "Funcional"

# 5.4 Actualizar -> 200
PUT {{base_url}}/api/clases/{{clase_id}}
Content-Type: application/json
{
  "nombre": "Funcional avanzado",
  "max_alumnos": 25
}
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Funcional avanzado"
jsonpath "$.max_alumnos" == 25


# =============================================================================
# 6. CLASE_HORARIO  (relación anidada clase -> horario)
# =============================================================================

# 6.1 Agregar horario a la clase -> 201, captura del id
POST {{base_url}}/api/clases/{{clase_id}}/horarios
Content-Type: application/json
{
  "dia_semana": "lunes",
  "hora_inicio": "19:00"
}
HTTP 201
[Captures]
horario_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.clase_id" == {{clase_id}}
jsonpath "$.dia_semana" == "lunes"
jsonpath "$.hora_inicio" startsWith "19:00"

# 6.2 Listar horarios de la clase
GET {{base_url}}/api/clases/{{clase_id}}/horarios
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1
jsonpath "$[0].id" == {{horario_id}}

# 6.3 Obtener horario por id dinámico
GET {{base_url}}/api/horarios/{{horario_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{horario_id}}
jsonpath "$.clase_id" == {{clase_id}}
jsonpath "$.dia_semana" == "lunes"


# =============================================================================
# 7. ALUMNO  (alta compuesta: usuario + alumno)
# =============================================================================

# 7.1 Crear alumno -> 201, captura del id
POST {{base_url}}/api/alumnos
Content-Type: application/json
{
  "email": "alumno@example.com",
  "contrasena": "secret123",
  "nombre": "Ana",
  "apellido": "Lopez",
  "fecha_inscripcion": "2026-10-01",
  "tipo_plan": "mensual"
}
HTTP 201
[Captures]
alumno_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.email" == "alumno@example.com"
jsonpath "$.rol" == "alumno"
jsonpath "$.fecha_inscripcion" == "2026-10-01"
jsonpath "$.tipo_plan" == "mensual"

# 7.2 Listar alumnos
GET {{base_url}}/api/alumnos
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1

# 7.3 Obtener alumno por id dinámico
GET {{base_url}}/api/alumnos/{{alumno_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{alumno_id}}
jsonpath "$.nombre" == "Ana"
jsonpath "$.apellido" == "Lopez"

# 7.4 Actualizar plan/vencimiento -> 200
PUT {{base_url}}/api/alumnos/{{alumno_id}}
Content-Type: application/json
{
  "fecha_vto": "2026-11-01",
  "tipo_plan": "trimestral"
}
HTTP 200
[Asserts]
jsonpath "$.tipo_plan" == "trimestral"
jsonpath "$.fecha_vto" == "2026-11-01"


# =============================================================================
# 8. INSCRIPCION  (relación anidada horario <-> alumno, valida cupo)
# =============================================================================

# 8.1 Inscribir alumno al horario -> 201
POST {{base_url}}/api/horarios/{{horario_id}}/inscripciones
Content-Type: application/json
{
  "alumno_id": {{alumno_id}}
}
HTTP 201
[Asserts]
jsonpath "$.clase_horario_id" == {{horario_id}}
jsonpath "$.alumno_id" == {{alumno_id}}
jsonpath "$.fecha_inscripcion" isString

# 8.2 Listar inscriptos del horario -> 200 y contiene al alumno
GET {{base_url}}/api/horarios/{{horario_id}}/inscripciones
HTTP 200
[Asserts]
jsonpath "$" isList
jsonpath "$" count >= 1
jsonpath "$[0].id" == {{alumno_id}}
jsonpath "$[0].nombre" == "Ana"


# =============================================================================
# 9. Casos de error y validaciones (edge cases)
# =============================================================================

# 9.1 GET de un recurso inexistente -> 404
GET {{base_url}}/api/ejercicios/999999
HTTP 404
[Asserts]
jsonpath "$.error" == "no encontrado"

# 9.2 POST con datos inválidos (nombre vacío) -> 400
POST {{base_url}}/api/ejercicios
Content-Type: application/json
{
  "nombre": ""
}
HTTP 400
[Asserts]
jsonpath "$.error" isString

# 9.3 POST con campo desconocido -> 400 (DisallowUnknownFields)
POST {{base_url}}/api/ejercicios
Content-Type: application/json
{
  "nombre": "x",
  "campo_raro": 1
}
HTTP 400
[Asserts]
jsonpath "$.error" == "JSON inválido"

# 9.4 GET con id no numérico -> 400
GET {{base_url}}/api/ejercicios/abc
HTTP 400
[Asserts]
jsonpath "$.error" isString

# 9.5 POST con email duplicado -> 409 (violación de unicidad)
POST {{base_url}}/api/profesores
Content-Type: application/json
{
  "email": "profe@example.com",
  "contrasena": "otra123",
  "nombre": "Otro",
  "apellido": "Profe"
}
HTTP 409


# =============================================================================
# 10. Eliminación (DELETE) y verificación posterior
#     Orden respetando claves foráneas: primero los hijos, luego los padres.
# =============================================================================

# 10.1 Eliminar inscripción -> 204
DELETE {{base_url}}/api/horarios/{{horario_id}}/inscripciones/{{alumno_id}}
HTTP 204

# 10.2 Quitar ejercicio de la rutina -> 204
DELETE {{base_url}}/api/rutinas/{{rutina_id}}/ejercicios/{{ejercicio_id}}
HTTP 204

# 10.3 Eliminar horario -> 204
DELETE {{base_url}}/api/horarios/{{horario_id}}
HTTP 204

# 10.4 Eliminar clase -> 204
DELETE {{base_url}}/api/clases/{{clase_id}}
HTTP 204

# 10.5 Eliminar rutina -> 204
DELETE {{base_url}}/api/rutinas/{{rutina_id}}
HTTP 204

# 10.6 Eliminar ejercicio -> 204
DELETE {{base_url}}/api/ejercicios/{{ejercicio_id}}
HTTP 204

# 10.7 Eliminar alumno -> 204
DELETE {{base_url}}/api/alumnos/{{alumno_id}}
HTTP 204

# 10.8 Eliminar profesor -> 204
DELETE {{base_url}}/api/profesores/{{profesor_id}}
HTTP 204

# 10.9 Consultar recursos ya eliminados -> 404
GET {{base_url}}/api/ejercicios/{{ejercicio_id}}
HTTP 404
[Asserts]
jsonpath "$.error" == "no encontrado"

GET {{base_url}}/api/alumnos/{{alumno_id}}
HTTP 404

GET {{base_url}}/api/profesores/{{profesor_id}}
HTTP 404

# 10.10 Volver a borrar un recurso ya eliminado -> 404
DELETE {{base_url}}/api/profesores/{{profesor_id}}
HTTP 404
