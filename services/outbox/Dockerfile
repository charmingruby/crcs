ARG GO_VERSION=1.27.0

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src/services/outbox

COPY pkg/ /src/pkg/
COPY services/outbox/ ./

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /src/bin/outbox ./cmd/outbox

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /src/bin/outbox /app/bin

ENTRYPOINT ["/app/bin"]
