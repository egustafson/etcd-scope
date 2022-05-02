#
BINS = etcd-scope

.PHONY: all clean etcd-scope

all: ${BINS} go.mod

etcd-scope: # $(shell find . -name \*.go)
	go build

clean:
	go clean --cache ./...
	-rm -f ${BINS}
