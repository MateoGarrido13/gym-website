-- Agrega el rol administrador. init.sh aplica cada archivo con su propio psql
-- (transacción aparte), así que el valor nuevo ya se puede usar en el seed.
-- IF NOT EXISTS lo hace idempotente para `make migrate` sobre una base existente.
ALTER TYPE enum_rol ADD VALUE IF NOT EXISTS 'administrador';
