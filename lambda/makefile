all: build zip
build:
	@GOOS=linux GOARCH=amd64 go build -o hello
zip:
	@zip handler.zip ./hello