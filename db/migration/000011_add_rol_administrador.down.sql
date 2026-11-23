-- PostgreSQL no permite quitar un valor de un enum. Para revertir hay que
-- borrar los usuarios administradores y recrear el tipo (o `make reset`).
DELETE FROM usuario WHERE rol = 'administrador';
