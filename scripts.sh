# consul
docker run -d -p 8500:8500 -p 8600:8600/udp --name=dev-consul consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client=0.0.0.0

# generate from proto file
protoc -I=api --go_out=. movie.proto
protoc -I=api --go_out=. --go-grpc_out=. movie.proto

# mysql
docker run --name movieexample_db -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=movieexample -p 3306:3306 -d mysql:latest
docker exec -i movieexample_db mysql movieexample -h localhost -P 3306 --protocol=tcp -uroot -ppassword < schema/schema.sql
docker exec -i movieexample_db mysql movieexample -h localhost -P 3306 --protocol=tcp -uroot -ppassword -e "SHOW tables"

# gRPC url
grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie", "user_id": "alex", "rating_value": 5}' localhost:8082 RatingService/PutRating
grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie"}' localhost:8082 RatingService/GetAggregatedRating

# Docker build
GOOS=linux go build -o main cmd/*.go
docker build -t metadata .
docker build -t rating .
docker build -t movie .

# Docker publish
docker tag metadata fuadaghazada/metadata:1.0.0
docker push fuadaghazada/metadata:1.0.0
docker tag metadata fuadaghazada/rating:1.0.0
docker push fuadaghazada/rating:1.0.0
docker tag metadata fuadaghazada/movie:1.0.0
docker push fuadaghazada/movie:1.0.0

# K8s
kubectl apply -f kubernetes-deployment.yml
kubectl get deployments
kubectl get pods
kubectl logs -f <POD_ID>
kubectl port-forward <POD_ID> 8081:8081