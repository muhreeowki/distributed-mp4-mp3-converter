build-auth:
	@cd auth-service && go build -o ../bin/auth-server

run-auth: build-auth
	@./bin/auth-server

build-gateway:
	@cd gateway-service && go build -o ../bin/gateway

run-gateway: build-gateway
	@./bin/gateway
