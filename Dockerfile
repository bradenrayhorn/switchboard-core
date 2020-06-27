FROM golang:1.14.4 as build

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest
COPY --from=build /app/switchboard-core /app/

EXPOSE 8080
EXPOSE 9001

ENTRYPOINT ["/app/switchboard-core"]
