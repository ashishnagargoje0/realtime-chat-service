FROM golang:1.24.4

WORKDIR /app
COPY . .
RUN go mod tidy

CMD ["go", "run", "cmd/server/main.go"]
