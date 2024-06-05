gengrt:
	protoc --go_out=./greet --go-grpc_out=./greet greet/greetpb/greet.proto

gencalc:
	protoc --go_out=./calculator --go-grpc_out=./calculator calculator/calculatorpb/calculator.proto

calser:
	go run calculator/calculator_server/server.go

calcli:
	go run calculator/calculator_client/client.go

grtser:
	go run greet/greet_server/server.go

grtcli:
	go run greet/greet_client/client.go

.PHONY: gengrt gencalc calser calcli grtser grtcli