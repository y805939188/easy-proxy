build_service:
	go build -o ./service/easy-proxy-service ./service/main.go
	go-bindata -o=./binary_service/service.go -pkg=binary_service ./service/easy-proxy-service
