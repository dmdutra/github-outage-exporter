FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /github-outage-exporter ./cmd/github-outage-exporter

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /github-outage-exporter /github-outage-exporter
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/github-outage-exporter"]
