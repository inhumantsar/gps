# Set the name and version of your application
APP_NAME=gps
APP_VERSION=latest

# Set the name and version of your Docker image
DOCKER_IMAGE=$(APP_NAME):$(APP_VERSION)

.PHONY: test debug

test:
	docker build -t $(DOCKER_IMAGE) -f test.Dockerfile .
	docker run $(DOCKER_IMAGE)

debug:
	docker build -t $(DOCKER_IMAGE) -f test.Dockerfile .
	docker run -it -v $(PWD):/app $(DOCKER_IMAGE) /bin/ash
