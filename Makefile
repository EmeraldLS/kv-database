db-server:
	go run server/main.go -p 3032

db-client:
	go run client/main.go -addr 3032

.PHONY: db-server db-client