IMAGE_NAME := go-rssag

.PHONY: build-run docker-build docker-run


build-run:
	go build -o bin/$(IMAGE_NAME) && ./$(IMAGE_NAME)

docker-build:
	docker build -f Dockerfile -t $(IMAGE_NAME):latest .

docker-run:
	docker run -p 8000:8000 $(IMAGE_NAME):latest