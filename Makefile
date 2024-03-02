run:
	go run cmd/app/main.go

register:
	go run cmd/register/main.go

deregister:
	rm db/wpp.db

build:
	go build -o app cmd/app/main.go && go build -o register cmd/register/main.go

build-armv6-app:
	GOOS=linux \
	GOARCH=arm \
	GOARM=6 \
	CGO_ENABLED=1 \
	CC=arm-linux-gnueabi-gcc \
	go build -o app cmd/app/main.go

build-armv6-register:
	GOOS=linux \
	GOARCH=arm \
	GOARM=6 \
	CGO_ENABLED=1 \
	CC=arm-linux-gnueabi-gcc \
	go build -o register cmd/register/main.go

build-armv6: build-armv6-app build-armv6-register