
build:
	go get
	go build

release: build
	git tag v0.2 && git push tags
	git-release release --user vsco --repo jsonconsul --tag v0.2

test: build
	go get golang.org/x/tools/cmd/cover
	go test -cover
