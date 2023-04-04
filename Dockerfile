FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 bubble
RUN go build -o BUPT-lib .

FROM alpine:latest AS tmp

RUN apk add --no-cache wget

RUN wget -O /build/BUPT-lib https://ghproxy.com/https://github.com/BUPTass/BUPT-lib/releases/latest/download/BUPT-lib

###################
# 接下来创建一个小镜像
###################
FROM alpine:latest

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=tmp /build/BUPT-lib /

# 需要运行的命令
ENTRYPOINT ["/BUPT-lib"]