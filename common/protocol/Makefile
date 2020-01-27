all: auth/auth.pb.go
.PHONY: all

auth/auth.pb.go: auth.proto
	protoc auth.proto --go_out=plugins=grpc:auth
