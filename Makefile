fmt:
	go fmt *.go
	# go fmt cmd/*.go
	go fmt sfomuseum/*.go

tools: 	
	go build -mod vendor -o bin/lookup cmd/lookup/main.go
