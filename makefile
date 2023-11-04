all: build vet lint

build:
	go build

vet:
	go vet

lint:
	golangci-lint run

release:
	ifndef VERSION
	$(error VERSION is not set)
	endif
	git tag $(VERSION)
	git push origin $(VERSION)
	GOPROXY=proxy.golang.org go list -m bandr.me/p/count@$(VERSION)
