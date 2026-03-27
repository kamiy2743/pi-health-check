FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY ./ ./
RUN go build -o ./health-check ./cmd/

FROM  debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/health-check /app/health-check
COPY ./urls /app/urls

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 1000 go
RUN useradd --uid 1000 --gid 1000 --no-create-home --shell /usr/sbin/nologin go
USER go:go

CMD ["/app/health-check"]
