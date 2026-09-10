FROM golang:1.26.5-alpine AS builder

ARG TARGET
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/service ./application/${TARGET}

FROM alpine:3.22

ARG TARGET
RUN apk add --no-cache ca-certificates tzdata wget curl \
    && addgroup -S microservice-kit \
    && adduser -S -G microservice-kit microservice-kit
WORKDIR /app
COPY --from=builder /out/service /app/service
COPY application/${TARGET}/conf.prod.yaml /app/conf.prod.yaml
COPY application/${TARGET}/zaplogger.prod.yaml /app/zaplogger.prod.yaml
COPY scripts/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chown -R microservice-kit:microservice-kit /app
USER microservice-kit
ENTRYPOINT ["/app/docker-entrypoint.sh"]
