api.pb.go api_grpc.pb.go: api.proto
	protoc -I. -I.. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--go-grpc_out=./gen  \
	--grpc-gateway_out=./gen   \
	--grpc-gateway_opt paths=source_relative  \
	--grpc-gateway_opt generate_unbound_methods=true  \
	--php_out=./gen  \
	--plugin=protoc-gen-grpc=$(which grpc_php_plugin) \
	--go_out=./gen api.proto
