# build stage
FROM golang:1.12 AS builder

WORKDIR /go/src/gitlab.com/vvelikodny/golang-microservices-test/query-client-service
COPY . /go/src/gitlab.com/vvelikodny/golang-microservices-test/query-client-service

RUN make query-client-service

# run stage
FROM alpine
COPY --from=builder /go/src/gitlab.com/vvelikodny/golang-microservices-test/bin/query-client-service /query-client-service

ENTRYPOINT ["/query-client-service"]
