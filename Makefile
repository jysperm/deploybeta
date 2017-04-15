OUTPUT=./releases
REPO=github.com/jysperm/deploying

all: release resouces

release:
	GOOS=darwin go build -o $(OUTPUT)/deploying-darwin-amd64
	GOOS=linux go build -o $(OUTPUT)/deploying-linux-amd64

resouces:
	cd frontend && gulp
	mkdir -p $(OUTPUT)/frontend/public
	cp -r frontend/public $(OUTPUT)/frontend/public

test:
	go test -v $(REPO)/lib/builder
	go test -v $(REPO)/lib/swarm
	go test -v $(REPO)/web/handlers
	go test -v $(REPO)/web/tests

clean:
	rm -r releases
