FROM golang:1.21 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -tags timetzdata -o app .

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y --no-install-recommends \
        tzdata \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*
ENV TZ Asia/Tokyo

WORKDIR /slack-event-notification-bot

COPY --from=builder /src/app .
COPY config.yaml .

ENTRYPOINT [ "/slack-event-notification-bot/app" ]
