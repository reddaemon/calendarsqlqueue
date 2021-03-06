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
	go build -o ./calendar
