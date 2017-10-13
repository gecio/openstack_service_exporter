.PHONY: all style vet lint code_style clean binary tarball version

BRANCH=$(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))
REVISION=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GITTAG=$(shell git describe --tags 2>/dev/null)
VERSIONFILE=$(shell tr -d "[:space:]" < VERSION)
VERSION=$(or $(GITTAG),$(VERSIONFILE))
GOARCH?=$(shell go env GOARCH)
GOOS?=$(shell go env GOOS)
CGO_ENABLED=0
BINARY=openstack_service_exporter

all: binary

version:
	@echo $(VERSION)

pkgs   = $(shell go list ./... | grep -v /vendor/)

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

vet:
	@echo ">> vetting code"
	@go vet $(pkgs)

lint:
	@echo ">> lint code"
	@golint -set_exit_status $(pkgs)

code_style: style vet lint

binary: $(BINARY)

$(BINARY):
	@echo ">> compile binary"
	@CGO_ENABLED=0 go build -ldflags "-X main.Branch=$(BRANCH) -X main.Version=$(VERSION) -X main.Revision=$(REVISION) -X main.BuildDate=$(BUILDDATE)" -o $(BINARY) .


tarball: binary
	@echo ">> creating tarball"
	@mkdir -p build/openstack_service_exporter-$(VERSION).$(GOOS)-$(GOARCH)
	@cp openstack_service_exporter build/openstack_service_exporter-$(VERSION).$(GOOS)-$(GOARCH)
	@cp LICENSE build/openstack_service_exporter-$(VERSION).$(GOOS)-$(GOARCH)
	@mkdir dist
	@tar czvf openstack_service_exporter-$(VERSION).$(GOOS)-$(GOARCH).tar.gz -C build openstack_service_exporter-$(VERSION).$(GOOS)-$(GOARCH)

docker:
	@echo ">> creating innovo/openstack_service_exporter:latest"
	@docker build -t innovo/openstack_service_exporter:latest .

clean:
	@echo ">> clean"
	@rm -fr ./dist/ ./build/ openstack_service_exporter openstack_service_exporter-*.tar.gz
