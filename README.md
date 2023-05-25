# annotations
Example of using grpc-gateway annotations for a simple health check:

```protobuf
service TaskManagementService {
    rpc Health(HealthRequest) returns (HealthResponse) {
        option (google.api.http) = {
            post: "/v1/example/echo"
            body: "*"
        };
    }
}
```

## Build & run
Run the following go commands to build and run the example:

```bash
go mod tidy
cd cmd/service
go run service.go
```
```log
# Output
2023/05/25 14:10:55 Starting service
2023/05/25 14:10:55 Listening http on: 127.0.0.1:8080
2023/05/25 14:10:55 Listening grpc on: 127.0.0.1:9090
```

## Call the service
```bash
curl -XPOST "localhost:8080/v1/example/echo" -H "content-type: text/json" -d "{}"
```
```bash
# Output
{"status":"healthy"}
```

## Regenerate protobuf code
The generate protobuf code is already checked into the repo,
but if you want to regenerate use the following commands at
the base of the repo:

```bash
rm -rf gen
buf generate proto --include-imports
```