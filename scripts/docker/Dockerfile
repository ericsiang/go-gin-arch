# Build stage
FROM golang:1.24-alpine AS builder

# 安裝必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 設置工作目錄
WORKDIR /app

# 複製 go mod 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/first_web_service ./cmd/first_web_service/main.go

# 構建遷移工具
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/migrate ./cmd/migrate/main.go

# Runtime stage
FROM alpine:latest

# 安裝 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata

# 設置時區為亞洲/台北
ENV TZ=Asia/Taipei

# 設置配置文件路徑
ENV CONFIG_PATH=/app/conf/

WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/bin/first_web_service /app/first_web_service
COPY --from=builder /app/bin/migrate /app/migrate

# 複製配置文件
COPY conf /app/conf

# 複製 log 目錄（如果需要）
RUN mkdir -p /app/log/info /app/log/error

# 暴露端口
EXPOSE 5001

# 運行應用
CMD ["/app/first_web_service"]
