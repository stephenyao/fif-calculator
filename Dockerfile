# ---------- 1) JS assets (build with esbuild) ----------
FROM node:20-alpine AS assets
WORKDIR /app

# Copy only what's needed for bundling to keep this stage cacheable/fast
# Adjust paths if your source JS lives elsewhere.
COPY static/src/ ./src/js/
# If you already have other static files (css/images), bring them too:
COPY static/ ./static/

# Bundle+minify your JS into /app/static/
# (keeps filename stable; for cache-busting see notes below)
RUN npx --yes esbuild src/js/fifCalculate.js \
    --bundle \
    --minify \
    --sourcemap \
    --outfile=static/fifCalculate.js


# ---------- 2) Go builder ----------
FROM golang:1.24.3 AS builder
WORKDIR /app

RUN apt-get update && \
    apt-get install -y gcc sqlite3 libsqlite3-dev && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the whole repo (or be selective if you want tighter caching)
COPY . .

# Generate templ
RUN templ generate

# Build your server
RUN go build -o app ./cmd/server


# ---------- 3) Runtime ----------
FROM debian:bookworm-slim
WORKDIR /app

# SQLite runtime libs + TLS certs
RUN apt-get update && apt-get install -y libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*

# App binary and required files
COPY --from=builder /app/app .
COPY --from=builder /app/views ./views
COPY --from=builder /app/trades.db ./trades.db

# 👉 Bring in the built static files from the assets stage
COPY --from=assets /app/static ./static

EXPOSE 8080
CMD ["./app"]