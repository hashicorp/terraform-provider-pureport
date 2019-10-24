SWEEP?="tf_test,tf-test"
TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'pureport/disabled')
TESTARGS?=$("-parallel=2","-timeout=120m")
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=pureport
PROVIDER_VERSION?=dev
GOOS?=darwin
GOARCH?=amd64
GOLINT_GOGC=5
GOLINT_JOBS=5

default: build

build: fmtcheck
	go install -ldflags="-X=github.com/terraform-providers/terraform-provider-pureport/version.ProviderVersion=$(PROVIDER_VERSION)"

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

install: plugin
	mkdir -p $(HOME)/.terraform.d/plugins/$(GOOS)_$(GOARCH)
	mv terraform-provider-pureport_$(PROVIDER_VERSION) $(HOME)/.terraform.d/plugins/$(GOOS)_$(GOARCH)

plugin: fmtcheck
	go build -ldflags="-X=github.com/terraform-providers/terraform-provider-pureport/version.ProviderVersion=$(PROVIDER_VERSION)" -o terraform-provider-pureport_$(PROVIDER_VERSION)

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -ldflags="-X=github.com/terraform-providers/terraform-provider-pureport/version.ProviderVersion=acc"

debugacc: fmtcheck
	TF_ACC=1 dlv test $(TEST) -- -test.v $(TESTARGS)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@GOGC=$(GOLINT_GOGC) golangci-lint run --concurrency $(GOLINT_JOBS) --deadline 4m ./$(PKG_NAME)
	@tfproviderlint -c 1 -S001 -S002 -S003 -S004 -S005 ./$(PKG_NAME)

tools:
	GO111MODULE=on go install github.com/bflad/tfproviderlint/cmd/tfproviderlint
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=on go install honnef.co/go/tools/cmd/staticcheck

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build sweep test testacc vet fmt fmtcheck errcheck lint tools test-compile website website-test

