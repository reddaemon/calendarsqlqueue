fmt:
	gofmt -w -s -d .
vet:
	 go vet ./...
lint:
	golint ./...
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
down:
	docker-compose -f docker/docker-compose/docker-compose.yml down
up-build:
	docker-compose -f docker/docker-compose/docker-compose.yml up --build -d
test:
	set -e ;\
	docker-compose -f docker/docker-compose/docker-compose-test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker/docker-compose/docker-compose-test.yml run integr_tests go test ./tests/integr || test_status_code=$$? ;\
	docker-compose -f docker/docker-compose/docker-compose-test.yml down ;\
	exit $$test_status_code ;\

