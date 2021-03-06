.DEFAULT_GOAL := all

examples = art landscape webtest conference nightvale poem scifi

.phony: fmt
fmt:
	go fmt . $(addprefix ./examples/, $(examples))

./bin/examples/%: *.go examples/**/*.go
	go build -o $@ ./examples/$(@F)

all: $(addprefix ./bin/examples/, $(examples))
	go test .