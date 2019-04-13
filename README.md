## Roles

Developer __Vitaly Velikodny__
  * @vvelikodny
  * [vvelikodny@gmail.com](mailto:vvelikodny@gmail.com)  

Requirements:
  * `protoc`
  * `go`
  * `docker`
  * `docker-compose`

## Deploy process (locally)


1.  Build Docker image

```bash
make
```

## Demo

```bash
curl -X POST http://localhost:8080/news -H "Content-Type: application/json" -d '{"Title": "News 1"}'
```

```bash
curl -X GET http://localhost:8080/news/1
```