.PHONY: all

OS = linux darwin windows
CURRENT_GOOS = $(shell go env GOOS)

pretty:
	gofmt -w *.go

sync:
	govendor sync

compile: pretty sync
	$(foreach var, $(OS), GOOS=$(var) GOARCH=amd64 CGO_ENABLED=0 go build -o build/jarpeace-$(var)-amd64 main.go;)

start: compile
	./build/jarpeace-$(CURRENT_GOOS)-amd64
	