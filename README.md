# Simple yet high performance load balancer

## Provide config
- Supports hot reload on config change
- Default config can be found at `./pkg/config/config.yaml`
```yaml
# load balancer health
health:
  interval: 2
  timeout: 5

## backends
backends:
  - name: backend1
    addr: localhost:8080
  - name: backend2
    addr: localhost:8081
  - name: backend3
    addr: localhost:8082
  - name: backend4
    addr: localhost:8083

# load balancer port
server:
  port: 8888
```

## ENV
- `CONFIG_PATH` = path to config file (optional)

## Running
```bash
go run main.go
```