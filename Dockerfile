FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ghbypass-client/ ./ghbypass-client/
COPY ghbypass-server/ ./ghbypass-server/
WORKDIR /app/ghbypass-server
RUN go mod download
RUN go build -o /app/server ./cmd/server
WORKDIR /app/ghbypass-client
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/client-linux ./cmd/client
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o /app/client-windows.exe ./cmd/client
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o /app/client-macos ./cmd/client

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /usr/local/bin/
COPY --from=builder /app/client-linux /usr/local/bin/client-linux
COPY --from=builder /app/client-windows.exe /usr/local/bin/client-windows.exe
COPY --from=builder /app/client-macos /usr/local/bin/client-macos
COPY --from=builder /app/server /usr/local/bin/ghbypass-server
ENTRYPOINT ["/usr/local/bin/ghbypass-server"]
