FROM golang:1.16-alpine AS build

COPY go.mod /go/src/httpserver/
COPY main.go /go/src/httpserver/
WORKDIR /go/src/httpserver/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/httpserver

FROM scratch
ENV VERSION=1.0
LABEL author="wangqi"
COPY --from=build /bin/httpserver /bin/httpserver
EXPOSE 80
ENTRYPOINT ["/bin/httpserver"]
