all: policy_proto server api client
server: server_proto protoset server_bin
api: api_proto api_bin
client: client_bin
api_proto:
	(cd api; protoc -I../policy/proto --proto_path=proto --go_out=plugins=grpc:proto --go_opt=paths=source_relative proto/api.proto)
server_proto:
	(cd server; protoc -I../policy/proto -I${GOPATH}/src/github.com/gogo/protobuf/gogoproto --proto_path=proto --gogo_out=plugins=grpc:proto --gogo_opt=paths=source_relative --include_source_info --include_imports --descriptor_set_out proto/genbyte/desc.protoset proto/server.proto)
policy_proto:
	(cd policy; protoc --proto_path=proto --go_out=plugins=grpc:proto --go_opt=paths=source_relative proto/policy.proto)
server_bin:
	go build -o build/server server/main.go server/dsc.go
api_bin:
	go build -o build/api api/main.go
client_bin:
	go build -o build/client client/main.go
protoset:
	(cd server/proto/genbyte; go run gen.go)
