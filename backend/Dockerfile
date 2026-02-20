# ビルドステージ
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

# 本番ステージ（distroless）
FROM gcr.io/distroless/static-debian12

COPY --from=builder /api /api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/api"]
