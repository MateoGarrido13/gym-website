#!/bin/sh
# Solo corre la PRIMERA vez que se crea el volumen (directorio de datos vacío).
# Orden: *.up.sql en orden numérico (enums → tablas) → funciones/procedures.
set -eu

for f in /opt/sql/migration/*.up.sql; do
    echo "Aplicando migración $f"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -f /opt/sql/migration/funciones.sql
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -f /opt/sql/migration/procedure.sql
