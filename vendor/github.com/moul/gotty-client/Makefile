COMMANDS :=	$(addprefix ./,$(wildcard cmd/*))
PACKAGES :=	 .
VERSION :=	$(shell cat .goxc.json | jq -c .PackageVersion | sed 's/"//g')
SOURCES :=	$(shell find . -name "*.go")
GOTTY_URL :=	http://localhost:8081/


.PHONY: build
build: $(notdir $(COMMANDS))


.PHONY: test
test: build
	./gotty-client $(GOTTY_URL)


.PHONY: build-docker
build-docker: contrib/docker/.docker-container-built


.PHONY: clean
clean:
	rm -f $(notdir $(COMMANDS))


.PHONY: install
install:
	go install $(COMMANDS)


.PHONY: release
release:
	goxc


.PHONY: build-docker
build-docker: contrib/docker/.docker-container-built


.PHONY: run-docker
run-docker: build-docker
	docker run -it --rm moul/gotty-client $(GOTTY_URL)


$(notdir $(COMMANDS)): $(SOURCES)
	gofmt -w $(PACKAGES) ./cmd/$@
	go test -i $(PACKAGES) ./cmd/$@
	go build -ldflags "-X main.VERSION=$(VERSION)" -o $@ ./cmd/$@
	./$@ --version


dist/latest/gotty-client_latest_linux_386: $(SOURCES)
	mkdir -p dist
	rm -f dist/latest
	(cd dist; ln -s $(VERSION) latest)
	goxc -bc="linux,386" xc
	cp dist/latest/gotty-client_$(VERSION)_linux_386 dist/latest/gotty-client_latest_linux_386


contrib/docker/.docker-container-built: dist/latest/gotty-client_latest_linux_386
	cp dist/latest/gotty-client_latest_linux_386 contrib/docker/gotty-client
	docker build -t moul/gotty-client:latest contrib/docker
	docker tag moul/gotty-client:latest moul/gotty-client:$(VERSION)
	docker run -it --rm moul/gotty-client --version
	docker inspect --type=image --format="{{ .Id }}" moul/gotty-client > $@.tmp
	mv $@.tmp $@
