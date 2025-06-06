FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    git gcc sqlite3 libsqlite3-dev postgresql-client netcat \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o main cmd/server/main.go

RUN chmod +x /app/start.sh
ENTRYPOINT ["/app/start.sh"]