GOBIN := $(shell go env GOPATH)/bin
SQLC := $(GOBIN)/sqlc

.PHONY: test sqlc build up down

#PATH REAL, SIN PHONY para que se chequee cada vez que se ejecuta el make
$(SQLC):
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
	@echo "SQLC version: $$($(SQLC) version) installed"

#utiliza la ruta absoluta para garantizar que siempre se encuentre sqlc sin importar el entorno, haciendo el proceso más confiable y portable.
sqlc: $(SQLC)
	@$(SQLC) generate

#Se compila pero no se usa en el test
build: sqlc
	@go build -o main main.go

up:
	@docker compose down -v
	@docker compose up -d --wait

down:
	@docker compose down -v

test: build up
	@go test ./test/... -v -count=1; \
	status=$$?;  \
	make down; \
	exit $$status