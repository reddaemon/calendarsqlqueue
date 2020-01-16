fmt:
	gofmt -w -s -d .
vet:
	go vet .
lint:
	golint .
tidy:
	go mod tidy
verify:
	go mod verify
precommit: fmt vet lint tidy verify
gen-proto:
	protoc protofiles/event.proto --go_out=plugins=grpc:.
build: gen-proto
	go build -o calendar_api cmd/api/*.go
	go build -o calendar_scheduler cmd/scheduler/*.go
	go build -o calendar_sender cmd/notifier/*.go
	go build -o calendar_client client/*.go
up:
	docker-compose -f docker/docker-compose/docker-compose.yml up -d

