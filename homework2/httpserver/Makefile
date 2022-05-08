export tag=v0.0.1
root:
	export ROOT=github.com/kenga01/cncamp/homework2/httpserver

build:
	echo "building httpserver library"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t kengaw/httpserver:${tag} .

push: release
	echo "pushing kengaw/httpserver"
	docker push kengaw/httpserver:${tag}
