build_arith:
	@go mod tidy
	@go build -o ./target/client ./actors/client/arith_client.go
	@go build -o ./target/server ./actors/server/arith_server.go
	
exe_client:
	@./target/client

exe_server:
	@./target/server

clean:
	@rm -rf target