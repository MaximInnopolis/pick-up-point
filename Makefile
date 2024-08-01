# build docker image
build:
	docker-compose build

up-all:
	docker-compose up -d postgres app tests zookeeper kafka1 kafka2 kafka3

down:
	docker-compose down

# run docker image
up-db:
	docker-compose up -d postgres

stop-db:
	docker-compose stop postgres

start-db:
	docker-compose start postgres

down-db:
	docker-compose down postgres


up-service:
	docker-compose up -d app --build

stop-service:
	docker-compose stop app

start-service:
	docker-compose start app

down-service:
	docker-compose down app


# Migration
migrate:
	./migration.sh up


# Mock generation
.PHONY: generate-mock
generate-mock:
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...

# Test
test:
	$(info running tests...)
	go test -v ./...


# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN:=$(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

ORDER_PROTO_PATH:="api/proto/order/v1"

# Установка всех необходимых зависимостей
.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest


.PHONY: generate
generate: .bin-deps
	mkdir -p pkg/${ORDER_PROTO_PATH}
	protoc -I api/proto \
		${ORDER_PROTO_PATH}/order.proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg/${ORDER_PROTO_PATH} --go_opt=paths=source_relative\
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg/${ORDER_PROTO_PATH} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ./pkg/api/proto/order/v1  --grpc-gateway_opt  paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=./pkg/api/proto/order/v1




