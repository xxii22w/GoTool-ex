proto:
	protoc -I=pb \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	pb/*.proto

run: 
	@go build -o bin/trim
	@./bin/trim