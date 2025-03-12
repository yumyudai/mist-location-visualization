#VER := $(shell git rev-parse HEAD | tr -d "\n")
VER := current
all: clean locapid mistpolld
clean:
	rm -rf out

container: locapid-container mistpolld-container

locapid:
	mkdir -p out
	go build -o out/locapid cmd/locapid/main.go

locapid-container:
	docker build -t locapid:$(VER) -f build/locapid/Dockerfile .
	docker tag locapid:$(VER) locapid:latest

mistpolld:
	mkdir -p out
	go build -o out/mistpolld cmd/mistpolld/main.go

mistpolld-container:
	docker build -t mistpolld:$(VER) -f build/mistpolld/Dockerfile .
	docker tag mistpolld:$(VER) mistpolld:latest
