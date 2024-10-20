# for debug use make SHELL="sh -x"

.DEFAULT_GOAL := build

GK_HOST = "0.0.0.0"
GK_PORT = "8080"
GK_HTTP_PORT = "8081"

PG_USER = "gophkeeper-user"
PG_PASSWORD = "gophkeeper-password"
PG_DB = "db-gophkeeper"
PG_HOST = "localhost"
PG_PORT = "5433"
PG_DATABASE_DSN = "postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable"
PG_IMAGE = "postgres:16.3-bookworm"
PG_DOCKER_CONTEINER_NAME = "gophkeeper-pg-16.3"

# Используем := чтобы переменная содержала значение на на момент определения этой переменной, см
# https://ftp.gnu.org/old-gnu/Manuals/make-3.79.1/html_chapter/make_6.html#SEC59
# TODO: так нужно задать все переменные в Makefile
BUILD_VERSION := 0.0.1
BUILD_DATE := $(shell date -u +"%Y-%m-%d %H:%M:%S:%N %Z")
BUILD_COMMIT := $(shell git rev-parse HEAD)

##--------------------------------------------------------------------
## PROTOBUF INSTALL
##--------------------------------------------------------------------
.PHONY:protobuf-install
protobuf-install:
	sudo apt install -y protobuf-compiler && \
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	# TODO: отразить в документации установку переменной окружения
	# export PATH="${PATH}:$(shell go env GOPATH)/bin"

##--------------------------------------------------------------------
## EASYJSON INSTALL
##--------------------------------------------------------------------
.PHONY:easyjson-install
easyjson-install:
	go get github.com/mailru/easyjson && \
	go install github.com/mailru/easyjson/...@latest

##--------------------------------------------------------------------
## PROTOBUF GENERATE
##--------------------------------------------------------------------
## TODO: Версии пакетов должны быть одни и теже в рамках Makefile.

GO_PATH := $(shell go env GOPATH)
PROTOBUF_PATH := "./proto/v1"
PROTOBUF_GEN_PATH := "./internal/adapters/api/grpc/gen/proto/v1"

.PHONY:protobuf-generate
protobuf-generate:
	mkdir -pv ./internal/adapters/api/grpc/gen/proto && \
	protoc \
		-I${GO_PATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1/third_party/googleapis \
		-I${GO_PATH}/pkg/mod/github.com/protocolbuffers/protobuf@v5.27.3+incompatible/src \
		-I. \
		-I${PROTOBUF_PATH} \
		--grpc-gateway_opt=Mitems.proto=. \
		--grpc-gateway_opt=Musers.proto=. \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_out=${PROTOBUF_GEN_PATH} \
		--go_opt=Mitems.proto=. \
		--go_opt=Musers.proto=. \
		--go_opt=paths=source_relative \
		--go_out=${PROTOBUF_GEN_PATH} \
		--go-grpc_opt=Mitems.proto=. \
		--go-grpc_opt=Musers.proto=. \
		--go-grpc_out=${PROTOBUF_GEN_PATH} \
		--go-grpc_opt=paths=source_relative \
		items.proto \
		users.proto

##--------------------------------------------------------------------
## OPENAPI2 INSTALL
##--------------------------------------------------------------------
## TODO: Указать конкретные версии пакетов для install?

.PHONY:openapi2-install
openapi2-install:
	go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

##--------------------------------------------------------------------
## OPENAPI2 GENERATE
##--------------------------------------------------------------------
## TODO: Версии пакетов должны быть одни и теже в рамках Makefile.

.PHONY:openapi2-generate
openapi2-generate: openapi2-install
	mkdir -pv ./third_party/OpenAPI && \
	protoc \
		-I${GO_PATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1/third_party/googleapis \
		-I${GO_PATH}/pkg/mod/github.com/protocolbuffers/protobuf@v5.27.3+incompatible/src \
		-I${PROTOBUF_PATH} \
		--openapiv2_opt=Mitems.proto=. \
		--openapiv2_opt=Musers.proto=. \
		--openapiv2_out=./third_party/OpenAPI \
		items.proto \
		users.proto

##--------------------------------------------------------------------
## BUILD, TESTS, RUN
##--------------------------------------------------------------------

GOLANG_LDFLAGS := -ldflags "-X 'main.buildVersion=${BUILD_VERSION}' \
                            -X 'main.buildDate=${BUILD_DATE}' \
                            -X 'main.buildCommit=${BUILD_COMMIT}'"

.PHONY:build
build:
	go build -C ./cmd/server/ -o server -buildvcs=false ${GOLANG_LDFLAGS}
	go build -C ./cmd/client/ -o client -buildvcs=false ${GOLANG_LDFLAGS}

.PHONY:godoc
godoc:
	go install golang.org/x/pkgsite/cmd/pkgsite@latest
	~/go/bin/pkgsite -open .

.PHONY:statictest
statictest:
	# statictest не переваривает имя пакета third_party
	go vet -vettool=$$(which statictest) ./internal/... ./cmd/...

.PHONY:test
test: build statictest 
	go test -v -race -count=1 ./...

.PHONY: cover
cover:
	mkdir -pv ./cover && \
	go test -v -coverpkg=./... -coverprofile=./cover/cover.profile.tmp ./... && \
	cat ./cover/cover.profile.tmp \
		| grep -v "model.go" \
		| grep -v "${PROTOBUF_PATH}" \
		> ./cover/cover.profile && \
	rm ./cover/cover.profile.tmp && \
	go tool cover -func ./cover/cover.profile && \
	go tool cover -html ./cover/cover.profile -o ./cover/cover.html

.PHONY: cover-clean
cover-clean:
	rm -v -f ./cover/cover.profile ./cover/cover.html && \
	rm -v -r ./cover

.PHONY:server-run-with-args
server-run-with-args: build db-up
	chmod +x ./cmd/server/server && \
	./cmd/server/server \
		-log-level debug \
		-address ${GK_HOST}:${GK_PORT} \
		-http-address ${GK_HOST}:${GK_HTTP_PORT} \
		-dsn ${PG_DATABASE_DSN}

.PHONY:client-run-with-args
client-run-with-args:
	chmod +x ./cmd/client/client && \
	./cmd/client/client \
		-log-level debug \
		-log-file client.log \
		-address ${GK_HOST}:${GK_PORT}

.PHONY:client-run-with-args-2
client-run-with-args-2:
	chmod +x ./cmd/client/client && \
	./cmd/client/client \
		-log-level debug \
		-log-file client2.log \
		-address ${GK_HOST}:${GK_PORT}
##--------------------------------------------------------------------
## DB POSTGRESQL
##--------------------------------------------------------------------

.PHONY: db-up
db-up:
	PG_USER=${PG_USER} \
	PG_PASSWORD=${PG_PASSWORD} \
	PG_DB=${PG_DB} \
	PG_HOST=${PG_HOST} \
	PG_PORT=${PG_PORT} \
	PG_DATABASE_DSN=${PG_DATABASE_DSN} \
	PG_IMAGE=${PG_IMAGE} \
	PG_DOCKER_CONTEINER_NAME=${PG_DOCKER_CONTEINER_NAME} \
	docker compose -f ./docker-compose.yml up -d postgres

.PHONY: db-down
db-down:
	PG_USER=${PG_USER} \
	PG_PASSWORD=${PG_PASSWORD} \
	PG_DB=${PG_DB} \
	PG_HOST=${PG_HOST} \
	PG_PORT=${PG_PORT} \
	PG_DATABASE_DSN=${PG_DATABASE_DSN} \
	PG_IMAGE=${PG_IMAGE} \
	PG_DOCKER_CONTEINER_NAME=${PG_DOCKER_CONTEINER_NAME} \
	docker compose -f ./docker-compose.yml down postgres

##--------------------------------------------------------------------
## GOLANGCI-LINT
##--------------------------------------------------------------------
GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.59.1 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	rm -rf ./golangci-lint
