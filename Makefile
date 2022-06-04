
gobuild:
	go build

run:
	./goapropos --basepath /usr/share/go-1.18/src --find args --all=false

rundebugtest:
	./goapropos --basepath /usr/share/go-1.18/src --find args --all=false --debug 12 --test=true

goclean:
	go clean

gofmt:
	go fmt

gorun:
	go run goapropos

