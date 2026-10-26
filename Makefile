GOBIN := $(shell go env GOPATH)/bin
SQLC := $(GOBIN)/sqlc
TEMPL := $(GOBIN)/templ
TEMPL_VERSION := v0.3.1070

.PHONY: test test-api sqlc templ generate-templ build up down reset run migrate

#PATH REAL, SIN PHONY para que se chequee cada vez que se ejecuta el make
$(SQLC):
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
	@echo "SQLC version: $$($(SQLC) version) installed"

#utiliza la ruta absoluta para garantizar que siempre se encuentre sqlc sin importar el entorno, haciendo el proceso más confiable y portable.
sqlc: $(SQLC)
	@$(SQLC) generate

# Genera las vistas server-side de templ.
$(TEMPL):
	@go install github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION)

templ: $(TEMPL)
	@$(TEMPL) generate

generate-templ: templ

# Compila el servidor y genera el código necesario antes de compilar las vistas.
build: sqlc templ
	@go build -o main main.go

# Postgres debe estar arriba (make up). No baja contenedores al salir.
run: build
	@./main

# No usa -v: el volumen postgres_data se conserva (datos + schema ya aplicados).
# El sleep cubre el reinicio de Postgres después de correr init.sh: --wait
# puede volver cuando el healthcheck pasa y la base todavía corta conexiones.
up:
	@docker compose up -d --wait
	@sleep 3

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
	@make up
	@go test ./test/... -v -count=1; \
	status=$$?;  \
	make reset; \
	exit $$status

# Schema limpio, servidor, requests.hurl y limpieza.
# La URL base se configura con la variable `base_url` dentro de requests.hurl.
test-api: build reset
	@make up
	@./main & echo $$! > .api.pid
	@sleep 2
	@hurl --test requests.hurl; \
	status=$$?; \
	kill $$(cat .api.pid) || true; \
	rm -f .api.pid; \
	make reset; \
	exit $$status