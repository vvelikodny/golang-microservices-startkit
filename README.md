## Roles

Developer __Vitaly Velikodny__
  * [@vvelikodny](https://github.com/vvelikodny)
  * [vvelikodny@gmail.com](mailto:vvelikodny@gmail.com)  

## Requirements:
  * `protoc`
  * `go`
  * `docker`
  * `docker-compose`

## Deploy process (locally)

Build services and run Docker containers

```bash
make
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
