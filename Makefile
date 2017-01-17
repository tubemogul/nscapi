PACKAGE=nscapi
FULLNAME=github.com/tubemogul/${PACKAGE}

.PHONY: all

all: test race
 
test: fmt lint vet
	go test $(GO_EXTRAFLAGS) -v -cover -covermode=count ./...

lint:
	golint $(GO_EXTRAFLAGS) -set_exit_status

fmt:
	@if [ -n "`gofmt -l .`" ]; then \
	 	printf >&2 'Some files are not in the gofmt format. Please fix.'; \
 		exit 1; \
	fi

vet:
	go tool vet -v *.go

race:
	go test $(GO_EXTRAFLAGS) -v -race ./...

bench:
	go test $(GO_EXTRAFLAGS) -v -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...
	go tool pprof -top -lines -nodecount=25 ${PACKAGE}.test cpu.prof
	go tool pprof -text -lines -nodecount=25 -alloc_space ${PACKAGE}.test mem.prof
	go tool pprof -text -lines -nodecount=25 -alloc_objects ${PACKAGE}.test mem.prof

gocov:
	gocov test | gocov report
	# gocov test >/tmp/gocovtest.json ; gocov annotate /tmp/gocovtest.json MyFunc

dep:
	go get -f -u -v ./...

build: dep test
	go clean -v ${FULLNAME}
	go build -v ${FULLNAME}

install: dep test
	go clean -v ${FULLNAME}
	go install ${FULLNAME}
