# Step 1: Modules caching
FROM golang:1.18.2-alpine3.14 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.18.2-alpine3.14 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/storage ./internal/storage/cmd


# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /bin/storage /storage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/storage"]