build:
	protoc -I. --proto_path=$GOPATH/src:. --micro_out=. --go_out=. proto/consignment/consignment.proto
	docker build -t shippy-consignment-service .

run:
	docker run -d --net="host" \
		-p 50052 \
		-e MICRO_SERVER_ADDRESS=:50052 \
		-e MICRO_REGISTRY=mdns \
		-e DISABLE_AUTH=true \
		shippy-consignment-service

