FROM golang:1.13 as build-env
ENV NAME "calendar_sender"
WORKDIR /opt/${NAME}
COPY ././go.mod .
COPY ././go.sum .
RUN go mod download

FROM build-env AS build
ENV NAME "calendar_sender"
WORKDIR /opt/${NAME}
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/${NAME} cmd/notifier/*.go

FROM alpine
ENV NAME "calendar_sender"
WORKDIR /opt/${NAME}
COPY --from=build /opt/${NAME}/bin/${NAME} ./${NAME}
RUN apk add --no-cache tzdata
CMD ./${NAME}