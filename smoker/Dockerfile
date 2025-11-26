FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o smoker-server ./cmd/server

RUN CGO_ENABLED=0 GOOS=linux go build -o smoker-client ./cmd/client

# Stage 2: Server Image ==============================================================
FROM alpine:latest AS server

WORKDIR /root/

COPY --from=builder /app/smoker-server .

EXPOSE 8080 8081 8082

CMD ["./smoker-server"]

# Stage 3: Client Image ==============================================================
FROM alpine:latest AS client

WORKDIR /root/

COPY --from=builder /app/smoker-client .

CMD ["./smoker-client"]
