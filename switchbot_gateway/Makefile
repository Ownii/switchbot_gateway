BINARY_NAME=switchbot

build:
	env GOOS=linux GOARCH=arm64 go build -buildvcs=false -o bin/${BINARY_NAME} github.com/Ownii/switchbot/cmd/switchbot

deploy:
	scp bin/${BINARY_NAME} martin.foerster@192.168.178.170:~/

build-deploy:
	make build
	make deploy