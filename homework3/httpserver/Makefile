export tag=v0.0.2
root:
	export ROOT=gitee.com/Peter1900/cncamptask/homework3/httpserver

build:
	echo "building httpserver library"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t cncamptask/httpserver:${tag} .
	docker tag cncamptask/httpserver:${tag} peter1900moon/httpserver:${tag}

push: release
	echo "pushing peter1900moon/httpserver"
	docker push peter1900moon/httpserver:${tag}
