#################################
# STEP 1 build executable binary
#################################
FROM golang:1.24-alpine AS builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create dbuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"


WORKDIR $GOPATH/ticketing-backend
COPY . .
COPY ./database/migrations /migrations

# Fetch dependencies.
RUN go get -d -v

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/ticketing-backend .

#######################################
# STEP 3 build a small image production
#######################################
FROM scratch AS production

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group



# Copy our static executable
COPY --from=builder /go/bin/ticketing-backend /ticketing-backend
COPY --from=builder /migrations /database/migrations

# Use an unprivileged user.
USER appuser:appuser

# Run the integrator binary.
EXPOSE 8080

ENTRYPOINT ["/ticketing-backend"]