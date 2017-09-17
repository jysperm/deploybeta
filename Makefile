OUTPUT=./releases
REPO=github.com/jysperm/deploying

export WORKDIR=$(shell pwd)

all: pack-tar

binaries:
	GOOS=linux go build -o $(OUTPUT)/deploying-linux-amd64

resources:
	cd frontend && gulp
	mkdir -p $(OUTPUT)/frontend $(OUTPUT)/lib/builder/runtimes/{golang,node}
	cp -r frontend/public $(OUTPUT)/frontend
	rm -r $(OUTPUT)/frontend/public/scripts
	cp lib/builder/runtimes/golang/Dockerfile.template $(OUTPUT)/lib/builder/runtimes/golang
	cp lib/builder/runtimes/node/Dockerfile.template $(OUTPUT)/lib/builder/runtimes/node

pack-tar: binaries resources
	cd $(OUTPUT) && tar --exclude *.tar.gz -zcvf deploying-linux-amd64.tar.gz *

test:
	go test -v $(REPO)/lib/runtimes
	go test -v $(REPO)/lib/builder
	go test -v $(REPO)/lib/models/app
	go test -v $(REPO)/lib/swarm
	go test -v $(REPO)/web/handlers
	go test -v $(REPO)/web/tests

clean:
	rm -r releases
