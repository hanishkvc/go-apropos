
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
	go run goapropos --test=true

gotest:
	go test
	go test -test.run Ma -test.v

gobenchmark:
	go test -test.bench .

samples:
	./goapropos fmt
	./goapropos --find numcpu
	./goapropos --findpkg fmt
	./goapropos --findcmt "type 'x'"
	./goapropos --findcmt "type flags for"
	./goapropos --findcmt "type flags"
	./goapropos --findcmt "uid and gid"
