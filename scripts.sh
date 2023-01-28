# consul
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0

# generate from proto file
protoc -I=api --go_out=. movie.proto
protoc -I=api --go_out=. --go-grpc_out=. movie.proto