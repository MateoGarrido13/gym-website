# gym-website

## Pagina web para gimnasios

Descripcion: Tendra el objetivo de digitalizar las rutinas, agilizar la actualizacion de las mismas y canalizar la inscrpcion de clases
e informacion en la aplicacion. A su vez, mejorar la cercania con el cliente y mejorar su experiencia.

Funcionalidades:
-Usuario para profesor- 
*Carga rutinas con ejercicios. 
*Asigna las rutinas a los alumnos.
*Gestionar clases ( pack de clases)
*Dashboard con metricas por pesona de su evolucion (pesos,concurrencia,ritmos)

-Usuario para alumno-
*Ve la rutina que le fue asignada
*Accede a los ejercicios,ve imágenes/videos. 
*Ver calendario e inscribirme a una clase nueva.


-ENTIDADES-

*Usuario: datos_personales [email, contrasena,nombre,apellido,telefono]

*Alumno
-fecha de inscripcion,fecha_vto,tipo_plan, nro_rutina

*Profesor : horario, especialidad

*Plan: duracion (cantidad de semanas ), ejercicios

*Ejercicio:nombre,descripcion

*Clase: dia,hora, cupos

-EJECUCION-

*Mediante consola, lanzando el comando:
*'go run main.go'
*Levantara el servidor en el puerto 8080 sirviendo el archivo 'index.html'

