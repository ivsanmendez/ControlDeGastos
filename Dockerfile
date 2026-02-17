# syntax=docker/dockerfile:1

# --- Build React frontend ---
FROM node:22-alpine AS web-builder

WORKDIR /web

COPY web/package*.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# --- Build Go API ---
FROM docker.io/library/golang:1.23-alpine AS api-builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# --- Runtime stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Copy Go API binary
COPY --from=api-builder /api /api

# Copy React build
COPY --from=web-builder /web/dist /web/dist

ENV STATIC_DIR=/web/dist

EXPOSE 8080

ENTRYPOINT ["/api"]