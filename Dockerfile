FROM golang:1.17.1-alpine3.13 as builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o /api

FROM scratch
COPY --from=builder /butcher /butcher
ADD config.json /config.json
EXPOSE 53
CMD [ "/butcher", "-c", "config.json", "0.0.0.0:53" ]
