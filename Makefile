gengrt:
	protoc --go_out=./greet --go-grpc_out=./greet greet/greetpb/greet.proto

gencalc:
	protoc --go_out=./calculator --go-grpc_out=./calculator calculator/calculatorpb/calculator.proto

.PHONY: gengrt gencalc