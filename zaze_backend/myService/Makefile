.PHONY: build clean deploy

build:
	dep ensure -v
	set GOOS=linux
	go build -ldflags="-s -w" -o bin/generate_static_site generate_static_site/main.go
	go build -ldflags="-s -w" -o bin/post_to_s3 post_to_s3/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose