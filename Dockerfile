FROM golang:1.18-alpine3.14 AS builder
WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -o ./appbin

FROM alpine:3.14.3
WORKDIR /app
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakaarta
COPY --from=builder /build/appbin ./
CMD ["./appbin"]