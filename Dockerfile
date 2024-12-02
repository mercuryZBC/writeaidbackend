FROM golang:1.23-alpine

# 设置工作目录
WORKDIR /app

# 复制项目文件
COPY . .

# 下载依赖并构建应用
RUN go mod tidy
RUN go build -o main ./cmd/app/

# 暴露端口
EXPOSE 8000

# 运行应用
CMD ["./main"]