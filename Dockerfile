FROM golang:1.10-alpine AS build
RUN apk add --no-cache git
WORKDIR /go/src/github.com/shaxbee/tmplserver
COPY . .
RUN go get -d ./...
RUN CGO_ENABLED=0 go install ./...

FROM scratch
WORKDIR /opt/tmplserver
COPY --from=build /go/bin/tmplserver .
CMD ["/opt/tmplserver/tmplserver"]