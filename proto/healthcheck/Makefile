healthcheck.pb.go healthcheck_grpc.pb.go: healthcheck.proto
	protoc -I. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go-grpc_out=. --go_out=. healthcheck.proto