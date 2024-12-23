FROM --platform=$TARGETPLATFORM golang:1.23-alpine AS builder
ENV CGO_ENABLED 0
# ENV GOPROXY https://goproxy.cn,direct
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -a -o adminapp ./cmd/admin/main.go

FROM --platform=$TARGETPLATFORM alpine
WORKDIR /app
COPY configs/ /app/configs/
COPY --from=builder /workspace/adminapp /app/
RUN apk update && apk add tzdata
ENV TZ Asia/Shanghai
ENV RUN_MODE test
EXPOSE 8081
ENTRYPOINT ["/app/adminapp"]