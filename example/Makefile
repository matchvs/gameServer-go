DOCKER_BUILD=./docker_build
BINARY = $(DOCKER_BUILD)/gameserver_go
 
.PHONY: all test image clean
all: build

test:
	
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY) 不能修改
build:
	rm -fr $(DOCKER_BUILD)
	mkdir -p $(DOCKER_BUILD)
	cp -r conf $(DOCKER_BUILD)
	cp Dockerfile $(DOCKER_BUILD)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY)

# image 内容不能修改
image: 
	sudo docker build -t $(IMAGE) $(DOCKER_BUILD)
	sudo docker tag $(IMAGE) $(HARBOR)/$(IMAGE)
	sudo docker push $(HARBOR)/$(IMAGE)
	sudo docker rmi $(IMAGE) $(HARBOR)/$(IMAGE)