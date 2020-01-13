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
imports:
	goimports -w
precommit: fmt vet lint tidy verify imports
gen-proto:
	protoc protofiles/event.proto --go_out=plugins=grpc:.
build:
	gen-proto
	go build main.go