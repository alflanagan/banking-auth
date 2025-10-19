# obviously, never do this anywhere but on a local dev machine
run:
	echo "Making run..."; \
	SERVER_ADDRESS=localhost SERVER_PORT=8081 DB_USER=root DB_PASSWD=admin DB_PORT=3306 DB_ADDR=127.0.0.1 DB_NAME=banking \
	go run main.go

build:
	go build -o banking-auth main.go
