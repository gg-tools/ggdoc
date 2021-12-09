# Swagger UI Server

## Usage

1. Run Program

```shell
go run main.go
```

Or use Docker:

```shell
docker run -d -p 80:80 liamylian/swagger-ui-server:latest
```

2. Visit WebPage`http://localhost`

## Configuration

```shell
# 1. Configure HTTP Port
export BIND_ADDR=:8000
# 2. Configure Static Files Directory
export HTML_ROOT=/root/apis/
# 3. Copy your API definition files to /root/apis/
cp myapi.json /root/apis/
# 4. Run program
go run main.go
# 5. Visit http://localhost:8000
```