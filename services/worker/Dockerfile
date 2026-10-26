ARG GO_VERSION=1.27.0

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin \
    ./cmd/worker


FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/bin /app/bin

ENTRYPOINT ["/app/bin"]
