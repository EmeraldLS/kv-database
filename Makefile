db-serve:
	go run internal/main.go -p 3032

connect-db:
	go run conn/main.go -addr 3032

.PHONY: db-server connect-db