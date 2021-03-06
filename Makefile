NAME = jsonconsul
VERSION = $(shell cat version)

all: deps build

deps:
	go get -d -v

build: deps
	@mkdir -p bin/
	go build -ldflags="-X main.Version $(VERSION)" -o bin/$(NAME) cmd/*.go

xcompile: deps
	gox -ldflags="-X main.Version $(VERSION)" -output="build/jsonconsul_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" github.com/vsco/jsonconsul/cmd

release: clean build test xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/tgz
	for f in $(FILES); do \
		(cd $(shell pwd)/build && tar -zcvf tgz/$$f.tar.gz $$f); \
		echo $$f; \
	done
	git tag $(VERSION) && git push origin --tags

clean:
	rm -rf bin
	rm -rf build
	
vet:
	go get golang.org/x/tools/cmd/vet
	go vet

test: build vet
	go get golang.org/x/tools/cmd/cover
	go test -cover
