all: generator build
generator-docker:
	export GO111MODULE="on" && go get &&  go test -v -run TestDocker_fetch github.com/sealstore/runtime
generator-containerd:
	export GO111MODULE="on" && go get &&  go test -v -run TestContainerd_fetch github.com/sealstore/runtime
generator:
	go get -u github.com/jteeuwen/go-bindata/...
	go-bindata -pkg command -o install/command/assert.go install/command/
build:
	export GO111MODULE="on" && go get && go build -o runtime
test:
	./runtime print
	./runtime print -d
