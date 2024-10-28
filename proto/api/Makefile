api.pb.go api_grpc.pb.go: api.proto
	protoc -I. -I../ --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go-grpc_out=. --go_out=. api.proto