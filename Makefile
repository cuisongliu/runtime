all: generator build
generator-docker:
	export GO111MODULE="on" && go get &&  go test -v -run TestDocker_fetch github.com/sealstore/runtime
generator-containerd:
	export GO111MODULE="on" && go get &&  go test -v -run TestContainerd_fetch github.com/sealstore/runtime
generator:
	wget https://github.com/cuisongliu/go-bindata/releases/download/v1.0/go-bindata
	chmod a+x go-bindata
	sudo mv go-bindata /usr/local/bin/go-bindata
	go-bindata -pkg command -o install/command/assert.go install/command/
build:
	export GO111MODULE="on" && go get && go build -o runtime
test:
	./runtime print
	./runtime print -d
