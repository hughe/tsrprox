all: tsrprox tsrprox.linux

tsrprox: tsrprox.go
	go build .

tsrprox.linux: tsrprox.go
	GOOS=linux GOARCH=amd64 go build .


