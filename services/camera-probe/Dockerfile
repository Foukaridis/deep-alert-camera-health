FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod ./
# RUN go mod download # (Simulating download)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /camera-probe .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /camera-probe /camera-probe
ENTRYPOINT ["/camera-probe"]
