FROM golang:1.10 AS build
WORKDIR /go/src/github.com/shaxbee/tmplserver
COPY . .
RUN go get -d ./...
RUN go install ./...

FROM gcr.io/distroless/base
EXPOSE 80
WORKDIR /opt/tmplserver
COPY --from=build /go/bin/tmplserver .
ENTRYPOINT ["/opt/tmplserver/tmplserver"]
