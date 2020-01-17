FROM golang:1.13 as build-env
ENV NAME "tests"
WORKDIR /opt/${NAME}
COPY ././go.mod .
COPY ././go.sum .
RUN go mod download

COPY . .

CMD ["go", "test"]
