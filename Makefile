server: 
	go run cmd/server/main.go

test: 
	go test ./...

bench:
	go test -bench=. -benchmem ./...

cover: 
	go test -cover -coverprofile=test/cover/cover.out ./...
	go tool cover -html test/cover/cover.out -o test/cover/cover.html
