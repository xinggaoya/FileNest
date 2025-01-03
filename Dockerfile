# 构建阶段
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux

# 安装基本工具
RUN apk add --no-cache git make

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main .

# 生产阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 设置时区为中国时区
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 创建存储目录
RUN mkdir -p /app/storage

# 设置权限
RUN chmod +x /app/main

# 暴露端口
EXPOSE 9040

# 启动应用
CMD ["./main"]