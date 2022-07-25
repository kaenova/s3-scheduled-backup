FROM golang:1.18-alpine3.14 AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./appbin

FROM alpine:3.14.3
WORKDIR /app
COPY --from=builder /build/appbin ./
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakaarta
CMD ["./appbin"]