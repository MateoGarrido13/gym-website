GOBIN := $(shell go env GOPATH)/bin
SQLC := $(GOBIN)/sqlc

.PHONY: test test-api sqlc build up down reset run migrate

#PATH REAL, SIN PHONY para que se chequee cada vez que se ejecuta el make
$(SQLC):
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
	@echo "SQLC version: $$($(SQLC) version) installed"

#utiliza la ruta absoluta para garantizar que siempre se encuentre sqlc sin importar el entorno, haciendo el proceso más confiable y portable.
sqlc: $(SQLC)
	@$(SQLC) generate

# Compila el binario del servidor (los tests de Go no lo usan; test sí lo compila).
build: sqlc
	@go build -o main main.go

# Postgres debe estar arriba (make up). No baja contenedores al salir.
run: build
	@./main

# No usa -v: el volumen postgres_data se conserva (datos + schema ya aplicados).
up:
	@docker compose up -d --wait

# Baja contenedores y deja el volumen intacto.
down:
	@docker compose down

# Único target que borra el volumen. Para tests o para rehacer el schema desde cero.
reset:
	@docker compose down -v

# Aplica UN .up.sql nuevo sobre la base ya existente. Ejemplo:
#   make migrate FILE=db/migration/000010_nueva_tabla.up.sql
migrate:
	@test -n "$(FILE)" || (echo "uso: make migrate FILE=db/migration/XXXX_nombre.up.sql" && exit 1)
	@docker compose exec -T postgres-db psql -U postgres -d tp2_db -v ON_ERROR_STOP=1 -f /opt/sql/migration/$$(basename $(FILE))

test: build reset
	@docker compose up -d --wait
	@go test ./test/... -v -count=1; \
	status=$$?;  \
	make reset; \
	exit $$status

# Prueba la API REST con Hurl. Requiere el servidor levantado en otra terminal
# (make up && make run) y Hurl instalado (https://hurl.dev).
# La URL base se configura con la variable `base_url` dentro de requests.hurl.
test-api:
	@hurl --test requests.hurl