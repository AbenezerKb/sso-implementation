FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /
ADD . .
RUN go build -o bin/sso /cmd/main.go

FROM alpine:3.16.0
WORKDIR /

COPY --from=builder /bin/sso .
COPY --from=builder /config/example_config.yaml /config/config.yaml
COPY --from=builder /internal/constant/query/schemas /internal/constant/query/schemas

EXPOSE 8000
ENTRYPOINT [ "./sso" ]
