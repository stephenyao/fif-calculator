FROM golang:1.24.3

WORKDIR /app

RUN apt-get update && \
    apt-get install -y gcc sqlite3 libsqlite3-dev

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate

RUN go build -o app ./cmd/server

FROM debian:bookworm-slim

# Install required runtime libraries for SQLite
RUN apt-get update && apt-get install -y libsqlite3-0 ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the compiled binary and needed files from the builder stage
COPY --from=0 /app/app .
COPY --from=0 /app/views ./views
COPY --from=0 /app/trades.db ./trades.db

EXPOSE 8080

CMD ["./app"]