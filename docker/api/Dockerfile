FROM golang:1.13 as build-env
ENV NAME "calendar_api"
WORKDIR /opt/${NAME}
COPY ././go.mod .
COPY ././go.sum .
RUN go mod download

FROM build-env AS build
ENV NAME "calendar_api"
WORKDIR /opt/${NAME}
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/${NAME} cmd/api/*.go

FROM alpine
ENV NAME "calendar_api"
WORKDIR /opt/${NAME}
COPY --from=build /opt/${NAME}/bin/${NAME} ./${NAME}
RUN apk add --no-cache tzdata
CMD GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info ./${NAME} --config config.yaml