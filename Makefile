
gobuild:
	go build

run:
	./goapropos --findpkg fmt print
	#./goapropos --basepath /usr/share/go-1.18/src --find args --all=false

rundebugtest:
	./goapropos --debug 12 --test args
	#./goapropos --basepath /usr/share/go-1.18/src --find args --all=false --debug 12 --test=true

createcache:
	./goapropos --autocache=false --createcache

allpkgssymbols:
	./goapropos --findpkg . --find .

goclean:
	go clean

gofmt:
	go fmt

gorun:
	go run goapropos --test=true

gotest:
	go test
	go test -test.run Ma -test.v

gotestpkgbasepath:
	./goapropos --autocache=false --createcache
	go test -test.run Pkg
	go test -test.run NONE -test.bench Pkg
	./goapropos --findpkg fmt print

gobenchmark:
	go test -test.run NONE -test.bench .

samples:
	./goapropos fmt
	./goapropos --find numcpu
	./goapropos --findpkg fmt
	./goapropos --findcmt "type 'x'"
	./goapropos --findcmt "type flags for"
	./goapropos --findcmt "type flags"
	./goapropos --findcmt "uid and gid"
