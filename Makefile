build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/main.exe ./cmd/*.go

build:
	go build -o build/main ./cmd/*.go