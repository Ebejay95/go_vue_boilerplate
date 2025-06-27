FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN if [ ! -f go.mod ]; then \
        go mod init myapp; \
    fi && \
    go get github.com/gin-gonic/gin@v1.9.1 && \
    go get github.com/lib/pq@v1.10.9 && \
    go get github.com/gin-contrib/cors@v1.4.0 && \
    go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]