FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest as runner 

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/predictions/thingNames.json /app/predictions/thingNames.json

CMD ["/app/main"]
