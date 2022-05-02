```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

sudo protoc --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=require_unimplemented_servers=false:. \
    --go-grpc_opt=paths=source_relative \
    api/v1/*.proto
```