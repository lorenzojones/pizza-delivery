# --- Build stage -------------------------------------------------------------
FROM golang:1.24-alpine AS build
WORKDIR /src

# Cache modules first.
COPY go.mod go.sum ./
RUN go mod download

# Build a static binary.
COPY . .
ARG VERSION=docker
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/server ./cmd/server

# --- Run stage ---------------------------------------------------------------
FROM alpine:3.20
# tzdata: for DASHBOARD_TIMEZONE (e.g. Europe/London).
# ca-certificates: for HTTPS calls to the Aircall API.
RUN apk add --no-cache tzdata ca-certificates wget

WORKDIR /app
COPY --from=build /out/server /app/server
# Bake in the example config and sample data so the container runs out of the
# box in demo mode. These are typically overridden by volume mounts in compose.
COPY config /app/config
COPY data /app/data

ENV PORT=8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://localhost:8080/api/health >/dev/null 2>&1 || exit 1

CMD ["/app/server"]
