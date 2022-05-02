build_service:
	go build -o ./service/easy-proxy-service ./service/main.go
	go-bindata -o=./binary_service/service.go -pkg=binary_service ./service/easy-proxy-service

build_cli:
	go build .

build_ubuntu:
	rm -rf /usr/local/bin/easy-proxy
	go build .
	ln ./easy-proxy /usr/local/bin/

build:
	make build_ubuntu