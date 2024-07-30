# for debug use make SHELL="sh -x"

.DEFAULT_GOAL := build

# Используем := чтобы переменная содержала значение на на момент определения этой переменной, см
# https://ftp.gnu.org/old-gnu/Manuals/make-3.79.1/html_chapter/make_6.html#SEC59
# TODO: так нужно задать все переменные в Makefile
BUILD_VERSION := 0.0.1
BUILD_DATE := $(shell date -u +"%Y-%m-%d %H:%M:%S:%N %Z")
BUILD_COMMIT := $(shell git rev-parse HEAD)

.PHONY:protobuf-install
protobuf-install:
	# from https://grpc.io/docs/protoc-installation/#:~:text=Linux%2C%20using%20apt%20or%20apt%2Dget
	sudo apt install -y protobuf-compiler
	# from https://practicum.yandex.ru/learn/go-advanced/courses/65ce3d44-da98-4684-9499-465ff6cc6c64/sprints/226895/topics/30311053-9716-4af0-9a23-f4fa0725f918/lessons/fa184729-fbbd-4a1c-ae11-4e12f66b7f64/#:~:text=%D0%A3%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0%20%D1%83%D1%82%D0%B8%D0%BB%D0%B8%D1%82%20gRPC
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="${PATH}:$(go env GOPATH)/bin"

PROTOBUF_PATH := "./internal/adapters/api/grpc/protobuf"

.PHONY:protobuf-generate
protobuf-generate:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		${PROTOBUF_PATH}/auth/model.proto


GOLANG_LDFLAGS := -ldflags "-X 'main.buildVersion=${BUILD_VERSION}' \
                            -X 'main.buildDate=${BUILD_DATE}' \
                            -X 'main.buildCommit=${BUILD_COMMIT}'"

.PHONY:build
build:
	go build -C ./cmd/server/ -o server -buildvcs=false ${GOLANG_LDFLAGS}

.PHONY:godoc
godoc:
	go install golang.org/x/pkgsite/cmd/pkgsite@latest
	~/go/bin/pkgsite -open .

.PHONY:statictest
statictest:
	go vet -vettool=$$(which statictest) ./...

.PHONY:test
test: build statictest
	go test -v -race -count=1 ./...

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
    golangci/golangci-lint:v1.55.2 \
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
