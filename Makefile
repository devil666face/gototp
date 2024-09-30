.DEFAULT_GOAL := help
PROJECT_BIN = $(shell pwd)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
GOBIN = go
PATH := $(PROJECT_BIN):$(PATH)
GOARCH = amd64
LINUX_LDFLAGS = -extldflags '-static' -w -s -buildid=
WINDOWS_LDFLAGS = -extldflags '-static' -w -s -buildid=
GCFLAGS = "all=-trimpath=$(shell pwd) -dwarf=false -l"
ASMFLAGS = "all=-trimpath=$(shell pwd)"
APP = gototp

build: build-linux build-windows .crop ## Build all

release: build-linux build-windows build-darwin .crop zip ## Build release

zip:
	cd $(PROJECT_BIN) && tar -cvzf $(PROJECT_BIN)/$(APP)_linux_amd64.tar.gz $(APP)
	cd $(PROJECT_BIN) && tar -cvzf $(PROJECT_BIN)/$(APP)_windows_amd64.tar.gz $(APP).exe
	cd $(PROJECT_BIN) && tar -cvzf $(PROJECT_BIN)/$(APP)_darwin_amd64.tar.gz $(APP)_darwin

docker: ## Build with docker
	docker compose up --build --force-recreate || docker-compose up --build --force-recreate


build-linux: ## Build for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) \
	  $(GOBIN) build -ldflags="$(LINUX_LDFLAGS)" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
	  -o $(PROJECT_BIN)/$(APP) cmd/gototp/main.go

build-windows: ## Build for windows
	CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) \
	  $(GOBIN) build -ldflags="$(WINDOWS_LDFLAGS)" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
	  -o $(PROJECT_BIN)/$(APP).exe cmd/gototp/main.go

build-darwin: ## Build for darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) \
	  $(GOBIN) build -ldflags="$(LINUX_LDFLAGS)" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
	  -o $(PROJECT_BIN)/$(APP)_darwin cmd/gototp/main.go

build-termux: ## Build for termux
	# https://dl.google.com/android/repository/android-ndk-r27-linux.zip
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 \
		CC=$(shell pwd)/android-ndk-r27/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android24-clang \
		CXX=$(shell pwd)/android-ndk-r27/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android24-clang \
		$(GOBIN) build -ldflags="-w -s" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
		-o $(PROJECT_BIN)/$(APP)_termux cmd/gototp/main.go
	
.crop:
	strip $(PROJECT_BIN)/$(APP)
	strip $(PROJECT_BIN)/$(APP).exe
	objcopy --strip-unneeded $(PROJECT_BIN)/$(APP)
	objcopy --strip-unneeded $(PROJECT_BIN)/$(APP).exe

dev:
	find . -name "*.go" | entr -r make build

help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
