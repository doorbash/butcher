FROM golang:1.17.6-alpine3.15 as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o /app

FROM alpine:3.15
RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true
COPY --from=builder /app /app
ADD config.json /config.json
EXPOSE 53
CMD [ "/app", "-c", "config.json", "0.0.0.0:53" ]
