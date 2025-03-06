running:
	CompileDaemon -build="go build -o ./cmd/main ./cmd" -command=./cmd/main

run:
	go run ./cmd/main.go