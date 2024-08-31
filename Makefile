server: 
	go run cmd/server/main.go

test: 
	go test ./...

cover: 
	go test -cover ./...
