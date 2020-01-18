[![Build Status](https://travis-ci.com/vvelikodny/golang-microservices-startkit.svg?branch=master)](https://travis-ci.com/vvelikodny/golang-microservices-startkit)

## Roles

Developer __Vitaly Velikodny__
  * [@vvelikodny](https://github.com/vvelikodny)
  * [vvelikodny@gmail.com](mailto:vvelikodny@gmail.com)  

## Requirements:
  * `go`
  * `docker`
  * `docker-compose`
  
  * `protoc` - to rebuilt protobuf

## Deploy process (locally)

Build services and run Docker containers

```bash
make run-env
```

## Demo

Add news

```bash
curl -X POST http://localhost:8080/news -H "Content-Type: application/json" -d '{"Title": "News 1"}'
```

Fetch news by `id`

```bash
curl -X GET http://localhost:8080/news/1
```
