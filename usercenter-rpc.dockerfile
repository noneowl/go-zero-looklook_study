FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
COPY . .
COPY app/usercenter/cmd/rpc/etc /app/etc
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/usercenter-rpc app/usercenter/cmd/rpc/usercenter-rpc.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/usercenter-rpc /app/usercenter-rpc
COPY --from=builder /app/etc /app/etc

CMD ["./usercenter-rpc", "-f", "etc/usercenter-rpc.yaml"]

