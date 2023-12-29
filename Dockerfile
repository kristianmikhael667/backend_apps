FROM golang:alpine3.18 as builder
RUN mkdir /app
COPY ./application /app
COPY ./.env /app/.env
WORKDIR /app
RUN go build -o backend_apps .

FROM alpine:3.18 
RUN apk update && apk add dumb-init
RUN mkdir /app
WORKDIR /app 
COPY --from=builder /app/.env /app
COPY --from=builder /app/backend_apps /app
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ./backend_apps