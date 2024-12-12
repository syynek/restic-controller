FROM golang:1.23-alpine AS builder

RUN apk update && \
    apk add alpine-sdk && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o bin/restic-controller .


FROM restic/restic:0.17.3

RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /app/bin/restic-controller .

ENTRYPOINT ["./restic-controller"]
