NAME = "jsonconsul"
VERSION = $(shell cat version)

all: deps build

deps:
	go get -d -v

build: deps
	@mkdir -p bin/
	go build -o bin/$(NAME)

xcompile: deps
	gox -output="build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"


release: build test xcompile
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

test: build
	go get golang.org/x/tools/cmd/cover
	go test -cover
