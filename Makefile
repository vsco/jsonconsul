release:
	git tag v0.2 && git push tags
	go build 
	git-release release --user vsco --repo jsonconsul --tag v0.2
