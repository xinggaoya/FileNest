FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制本地文件到容器中
COPY . /app

# 设置Go 下载镜像源
ENV GOPROXY https://goproxy.cn

# 编译程序
RUN go build -o main .

FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制编译好的程序到容器中
COPY --from=builder /app/main /app/main

# 设定时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 需要暴露的端口号
EXPOSE 9040

# 设置可执行文件权限
RUN chmod +x /app/main

# 启动命令
CMD ["/app/main"]