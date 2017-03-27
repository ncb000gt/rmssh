build: deps
	go build

deps:
	go get golang.org/x/crypto/ssh
	go get golang.org/x/crypto/ssh/terminal
	go get github.com/fatih/color
