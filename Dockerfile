FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    git gcc sqlite3 libsqlite3-dev postgresql-client netcat \
  && rm -rf /var/lib/apt/lists/*
ARG DB_USER
ARG DB_PASSWORD
ARG DB_DATABASE
ARG DB_SSLMODE
ARG ENV
ARG JWT_SECRET
ARG EXP_TIME

# Set runtime environment variables
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_DATABASE=${DB_DATABASE}
ENV DB_SSLMODE=${DB_SSLMODE}
ENV ENV=${ENV}
ENV JWT_SECRET=${JWT_SECRET}
ENV EXP_TIME=${EXP_TIME}

RUN ECHO EXP_TIME
RUN ECHO JWT_SECRET

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o main cmd/server/main.go

RUN chmod +x /app/start.sh
ENTRYPOINT ["/app/start.sh"]