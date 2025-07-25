# s6g-order

## GRPC
### install grpc server
```bash
go get -u google.golang.org/grpc
```
### install protoc compiler
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

## SQLC
### install sqlc in global
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```


### sqlc compile
```bash
sqlc generate
```
o