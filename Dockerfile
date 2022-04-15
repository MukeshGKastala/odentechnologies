FROM golang:1.16-alpine AS build

WORKDIR /app

RUN apk add --no-cache make

COPY go.mod ./
COPY go.sum ./
RUN --mount=type=ssh go mod download
RUN --mount=type=ssh go mod verify

COPY . .
RUN make build

RUN adduser -D -g '' -s /bin/false -h /odentechnologies odentechnologies

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/bin/metrics /bin/metrics
COPY --from=build /etc/passwd /etc/passwd

USER odentechnologies

ENTRYPOINT ["/bin/metrics"]