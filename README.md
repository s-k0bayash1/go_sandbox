# go_sandbox
### Run on docker
```sh
docker build -t go-sandbox-1.18:latest .
docker run --rm -v $PWD:/app go-sandbox-1.18:latest go run generics/main.go
```
