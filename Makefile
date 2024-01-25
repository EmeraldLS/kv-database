db-server:
	go run server/main.go 

db-client:
	go run client/main.go 

.PHONY: db-server db-client