FROM golang:1.23.4-bullseye AS builder

WORKDIR /app
COPY . .

# we need to statically link for sqlite3
RUN CGO_ENABLED=1 go build \
  -tags "osusergo netgo static_build" \
  -ldflags "-s -w -linkmode external -extldflags '-static'" \
  -o /groxy main.go

FROM scratch

COPY --from=builder /groxy /groxy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set TERM environment variable for true color
ENV TERM=xterm-truecolor

ENTRYPOINT ["/groxy"]

EXPOSE 8080
