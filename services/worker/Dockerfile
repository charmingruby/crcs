ARG GO_VERSION=1.27.0

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src

COPY pkg/ ./pkg/
COPY services/ ./services/

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /src/bin/worker \
    ./services/worker/cmd/worker


FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /src/bin/worker /app/bin

ENTRYPOINT ["/app/bin"]
