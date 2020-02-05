all: dirs auth/auth.pb.go sensors/sensors.pb.go
.PHONY: all dirs

dirs: auth sensors

auth:
	mkdir auth

auth/auth.pb.go: auth.proto
	protoc auth.proto --go_out=plugins=grpc:auth

sensors:
	mkdir sensors

sensors/sensors.pb.go: sensors.proto
	protoc sensors.proto --go_out=plugins=grpc:sensors
