# Multi-stage: build with Go, run on minimal Alpine
FROM golang:1.26-alpine AS build
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILDDATE=unknown
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w \
	-X dcctl/cmd/dcctl/root.buildVersion=$(VERSION) \
	-X dcctl/cmd/dcctl/root.buildCommit=$(COMMIT) \
	-X dcctl/cmd/dcctl/root.buildDate=$(BUILDDATE)" \
	-o /dcctl ./cmd/dcctl

FROM alpine:3.19
LABEL org.opencontainers.image.title="dcctl"
LABEL org.opencontainers.image.description="Docker Compose Control - manage Compose stacks per environment"
LABEL org.opencontainers.image.source="https://github.com/hrodrig/dcctl"
RUN apk --no-cache add ca-certificates
RUN adduser -D -g "" dcctl
COPY --from=build /dcctl /home/dcctl/dcctl
RUN chown dcctl:dcctl /home/dcctl/dcctl
USER dcctl
WORKDIR /home/dcctl
ENTRYPOINT ["/home/dcctl/dcctl"]
