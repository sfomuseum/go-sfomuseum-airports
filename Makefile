cli:
	go build -mod vendor -o bin/lookup cmd/lookup/main.go
	go build -mod vendor -o bin/build-sfomuseum-data cmd/build-sfomuseum-data/main.go

rebuild:
	go build -mod vendor -o bin/build-sfomuseum-data cmd/build-sfomuseum-data/main.go
	bin/build-sfomuseum-data
	go build -mod vendor -o bin/lookup cmd/lookup/main.go
