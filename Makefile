build:
	go build -o bin/main main.go

run:
	go run main.go

migrate:
	cd database/
	go run gendb.go

buildweb:
	cd web/
	yarn run build

purge:
	rm -rf ./bin
	rm -rf ./web/build

install:
	go mod download
	go build -o bin/fiber main.go