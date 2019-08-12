.PHONY: build clean test default

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
BINARY_NAME=gitlab-status
IMAGE_NAME=$(BINARY_NAME)
BINARY_MAC=$(BINARY_NAME)-darwin
BINARY_WINDOWS=$(BINARY_NAME).exe
ZIP_LINUX=$(BINARY_NAME)-amd64-linux.zip
ZIP_WINDOWS=$(BINARY_NAME)-amd64-win.zip
ZIP_MAC=$(BINARY_NAME)-amd64-darwin.zip

all: test build
default: build
build:
	@echo "building ${BINARY_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X github.com/max-wittig/gitlab-status/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/max-wittig/gitlab-status/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_NAME}
run: build
	bin/$(BINARY_NAME) -daemon=false -config config.example.yml
run-daemon: build
	bin/$(BINARY_NAME) -config config.example.yml
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/max-wittig/gitlab-status/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/max-wittig/gitlab-status/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_NAME}
	cd bin && zip -r $(ZIP_LINUX) $(BINARY_NAME)
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/max-wittig/gitlab-status/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/max-wittig/gitlab-status/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_WINDOWS}
	cd bin && zip -r $(ZIP_WINDOWS) $(BINARY_WINDOWS)
build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X github.com/max-wittig/gitlab-status/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/max-wittig/gitlab-status/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_MAC}
	cd bin && zip -r $(ZIP_MAC) $(BINARY_MAC)
package-docker:
	@echo "building image ${BINARY_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(IMAGE_NAME):latest .
build-all: build-linux build-windows build-mac
clean:
	@test ! -e bin/${BINARY_NAME} || rm -rf bin/
test:
	go test ./...
