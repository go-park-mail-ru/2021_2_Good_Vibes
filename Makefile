.PHONY: build
build:
	go build -o ./build/api/main ./cmd/api
	go build -o ./build/order/main ./cmd/order
	go build -o ./build/basket/main ./cmd/basket
	go build -o ./build/auth/main ./cmd/auth
