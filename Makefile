gengrt:
	protoc --go_out=./greet --go-grpc_out=./greet greet/greetpb/greet.proto

.PHONY: gengrt