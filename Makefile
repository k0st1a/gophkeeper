# for debug use make SHELL="sh -x"

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
