.phony: fmt
fmt:
	go fmt . && go fmt ./examples/art && go fmt ./examples/landscape

./bin/examples/%: *.go examples/**/*.go
	go build -o $@ ./examples/$(@F)

all: ./bin/examples/art ./bin/examples/landscape