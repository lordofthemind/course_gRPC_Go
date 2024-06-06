gengrt:
	protoc --go_out=./greet --go-grpc_out=./greet greet/greetpb/greet.proto

gencalc:
	protoc --go_out=./calculator --go-grpc_out=./calculator calculator/calculatorpb/calculator.proto

genblg:
	protoc --go_out=./blog --go-grpc_out=./blog blog/blogpb/blog.proto

calser:
	go run calculator/calculator_server/server.go

calcli:
	go run calculator/calculator_client/client.go

grtser:
	go run greet/greet_server/server.go

grtcli:
	go run greet/greet_client/client.go

blgser:
	go run blog/blog_server/server.go

blgcli:
	go run blog/blog_client/client.go

.PHONY: gengrt gencalc calser calcli grtser grtcli genblg blgser blgcli