.PHONY: gen
gen:
	protoc --proto_path=api-v1-library/api/proto \
	--grpc-gateway_out=logtostderr=true:api-v1-library/api/proto \
    --proto_path=third_party \
    --go_out=plugins=grpc:api-v1-library/api/proto service.proto

	protoc --proto_path=svc-books/api/proto \
    --proto_path=third_party \
    --go_out=plugins=grpc:svc-books/api/proto service.proto

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ./builds/gw-library/gw-library -i ./api-v1-library/cmd/gw-library/*.go
	GOOS=linux GOARCH=amd64 go build -o ./builds/svc-library/svc-library -i ./api-v1-library/cmd/svc-library/*.go
	GOOS=linux GOARCH=amd64 go build -o ./builds/svc-books/svc-books -i ./svc-books/cmd/svc-books/*.go

.PHONY: image
image:
	docker build -t gw-library ./builds/gw-library
	docker build -t svc-library ./builds/svc-library
	docker build -t svc-books ./builds/svc-books
	docker build -t db-library ./db

.PHONY: up
up:
	docker-compose up

# clean
.PHONY: clean
clean:
	docker images -q -f "dangling=true" | xargs -I {} docker rmi {}
	@docker rmi -f "gw-library"
	@docker rmi -f "svc-library"
	@docker rmi -f "svc-books"
	@docker rmi -f "db-library"

.PHONY: all
all: gen build image up