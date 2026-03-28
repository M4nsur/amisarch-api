FROM golang:1.26.1-alpine3.23 AS build

WORKDIR /api


COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /myapp ./cmd/api

FROM alpine:3.22 AS run

COPY --from=build /myapp /myapp

EXPOSE 8080

CMD ["/myapp"]
