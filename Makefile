OS = darwin freebsd linux openbsd windows
ARCHS = 386 arm amd64 arm64

all: build release

deps:
	go get -d -v -t ./...

build: deps
	go build

install: deps
    go install

test: deps
	go test ./...

release: clean deps
	@for arch in $(ARCHS);\
	do \
		for os in $(OS);\
		do \
			echo "Building for $$os-$$arch"; \
			mkdir -p build/http-server-$$os-$$arch/; \
			GOOS=$$os GOARCH=$$arch go build -o build/http-server-$$os-$$arch/http-server; \
			tar cz -C build -f build/http-server-$$os-$$arch.tar.gz http-server-$$os-$$arch; \
		done \
	done

clean:
	rm -rf build
	rm -f http-server
