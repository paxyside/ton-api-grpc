FROM golang:1.24.2 AS build

LABEL builder="go build"

WORKDIR /build

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
COPY .golangci.yaml ./

RUN go mod download && go mod verify

COPY . .

RUN go tool golangci-lint run --fix
RUN go test -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w" -o /usr/bin/server ./cmd/server/main.go

FROM alpine:3.12

RUN apk update && \
    apk add --no-cache \
        ca-certificates \
        curl \
        tzdata \
    && rm -rf -- /var/cache/apk/*

WORKDIR /app
COPY --from=build /usr/bin/server .

HEALTHCHECK --interval=20s --timeout=5s --retries=5 --start-period=30s \
    CMD grpcurl -fsS -m5 -A'docker-healthcheck' http://127.0.0.1/tonnode.TonNodeService/Ping || exit 1

CMD ["./server"]
