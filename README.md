# GG Doc

- [x] OpenAPI Specification 3.0
    - [x] Swagger
    - [x] RedDoc
- [x] Markdown
    - [x] Highlight
    - [x] PlantUML
    - [ ] Mermaid
    - [ ] PDF Print
- [x] Text

## Usage

1. Run Program

```shell
go run main.go
```

Or use Docker:

```shell
docker run -d --restart always --name apidoc \
  -p 80:80 \
  -v ./docs:/root/docs \
  liamylian/ggdoc:latest
```

2. Visit WebPage`http://localhost`

## Configuration

```shell
# 1. Configure HTTP Port
export BIND_ADDR=:8000
# 2. Configure Static Files Directory
export DOCS_ROOT=/root/docs/
# 3. Copy your API definition files to /root/apis/
cp myapi.json /root/docs/
# 4. Run program
go run main.go
# 5. Visit http://localhost:8000
```