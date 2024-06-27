FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest as runner 

WORKDIR /app

COPY --from=builder /app/main /app/main

COPY --from=builder /app/predictions/thingNames.json /app/predictions/thingNames.json

COPY --from=builder /app/tracking/example_track_long_acc.csv /app/tracking/example_track_long_acc.csv
COPY --from=builder /app/tracking/example_track_long_gyro.csv /app/tracking/example_track_long_gyro.csv
COPY --from=builder /app/tracking/example_track_long_mag.csv /app/tracking/example_track_long_mag.csv
COPY --from=builder /app/tracking/example_track_long_gps.csv /app/tracking/example_track_long_gps.csv
COPY --from=builder /app/tracking/example_track_long_acc.csv.gz /app/tracking/example_track_long_acc.csv.gz
COPY --from=builder /app/tracking/example_track_long_gyro.csv.gz /app/tracking/example_track_long_gyro.csv.gz
COPY --from=builder /app/tracking/example_track_long_mag.csv.gz /app/tracking/example_track_long_mag.csv.gz
COPY --from=builder /app/tracking/example_track_long_gps.csv.gz /app/tracking/example_track_long_gps.csv.gz
COPY --from=builder /app/tracking/example_track_long.json.gz /app/tracking/example_track_long.json.gz
COPY --from=builder /app/tracking/example_track_long.json /app/tracking/example_track_long.json

COPY --from=builder /app/tracking/example_track_short_acc.csv /app/tracking/example_track_short_acc.csv
COPY --from=builder /app/tracking/example_track_short_gyro.csv /app/tracking/example_track_short_gyro.csv
COPY --from=builder /app/tracking/example_track_short_mag.csv /app/tracking/example_track_short_mag.csv
COPY --from=builder /app/tracking/example_track_short_gps.csv /app/tracking/example_track_short_gps.csv
COPY --from=builder /app/tracking/example_track_short_acc.csv.gz /app/tracking/example_track_short_acc.csv.gz
COPY --from=builder /app/tracking/example_track_short_gyro.csv.gz /app/tracking/example_track_short_gyro.csv.gz
COPY --from=builder /app/tracking/example_track_short_mag.csv.gz /app/tracking/example_track_short_mag.csv.gz
COPY --from=builder /app/tracking/example_track_short_gps.csv.gz /app/tracking/example_track_short_gps.csv.gz
COPY --from=builder /app/tracking/example_track_short.json.gz /app/tracking/example_track_short.json.gz
COPY --from=builder /app/tracking/example_track_short.json /app/tracking/example_track_short.json


CMD ["/app/main"]
