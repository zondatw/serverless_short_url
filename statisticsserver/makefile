GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
DOCKERHUB = zondayang
NAMESPACE = statistics_server
BINARY_NAME = statisticsserver

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

docker:
	@docker build --tag docker.io/$(DOCKERHUB)/$(NAMESPACE):latest .
	@docker push docker.io/$(DOCKERHUB)/$(NAMESPACE):latest

.PHONY: clean
clean:
	$(GOCLEAN)
	rm $(BINARY_NAME)