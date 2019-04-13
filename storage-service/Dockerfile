# build stage
FROM golang:1.12 AS builder

WORKDIR /go/src/gitlab.com/vvelikodny/golang-microservices-test/storage-service
COPY . /go/src/gitlab.com/vvelikodny/golang-microservices-test/storage-service

RUN make storage-service

# run stage
FROM alpine
COPY --from=builder /go/src/gitlab.com/vvelikodny/golang-microservices-test/storage-service/storage-service /storage-service

ENTRYPOINT ["/storage-service"]
