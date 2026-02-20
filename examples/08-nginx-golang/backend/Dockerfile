FROM golang:1.21-alpine AS builder
WORKDIR /code
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /backend main.go

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
COPY --from=builder /backend /backend
EXPOSE 80
CMD ["/backend"]
