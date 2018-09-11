.PHONY: all

pretty:
	gofmt -w *.go

sync:
	govendor sync

compile: pretty sync
	go build -o jarpeace main.go

start: compile
	./jarpeace