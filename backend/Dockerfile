FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /src

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM docker.io/library/alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/server /app/server
COPY memory /app/memory
COPY ragDocs /app/ragDocs

ENV SERVER_ADDR=:8081 \
    KNOWLEDGE_MEMORY_DIR=/app/memory

EXPOSE 8081

CMD ["/app/server"]
