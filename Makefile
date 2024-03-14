# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./...

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix ./...

##@ Build

.PHONY: build
build: fmt vet ## Build the code generator.
	go build -o $(LOCALBIN)/cwc ./main.go

.PHONY: run
run: fmt vet check-env-vars ## Run the example app.
	go run main.go

.PHONY: generate
generate: lang-gen ## Generate code.
	go generate ./...

##@ Pre Deployment

.PHONY: security-scan
security-scan: gosec ## Run security scan on the codebase.
	$(GOSEC) ./...

.PHONY: dependency-scan
dependency-scan: govulncheck ## Run dependency scan on the codebase.
	$(GOVULNCHECK) ./...

##@ Deployment


##@ Dependencies

.PHONY: check-env-vars
check-env-vars:
	@missing_vars=""; \
	for var in AOAI_API_KEY AOAI_ENDPOINT AOAI_API_VERSION AOAI_MODEL_DEPLOYMENT; do \
			if [ -z "$${!var}" ]; then \
					missing_vars="$$missing_vars $$var"; \
			fi \
	done; \
	if [ -n "$$missing_vars" ]; then \
			echo "Error: the following env vars are not set:$$missing_vars"; \
			exit 1; \
	fi

.PHONY: lang-gen
lang-gen:
	cd generator && go build -o $(LOCALBIN)/lang-gen lang-gen.go

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
TOOLKIT_TOOLS_GEN = $(LOCALBIN)/toolkit-tools-gen-$(TOOLKIT_TOOLS_GEN_VERSION)
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
GOSEC = $(LOCALBIN)/gosec-$(GOSEC_VERSION)
GOVULNCHECK = $(LOCALBIN)/govulncheck-$(GOVULNCHECK_VERSION)

## Tool Versions
TOOLKIT_TOOLS_GEN_VERSION ?= latest
GOLANGCI_LINT_VERSION ?= v1.54
GOSEC_VERSION ?= latest
GOVULNCHECK_VERSION ?= latest


.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

.PHONY: gosec
gosec: $(GOSEC) ## Download golangci-lint locally if necessary.
$(GOSEC): $(LOCALBIN)
	$(call go-install-tool,$(GOSEC),github.com/securego/gosec/v2/cmd/gosec,$(GOSEC_VERSION))

.PHONY: govulncheck
govulncheck: $(GOVULNCHECK) ## Download golangci-lint locally if necessary.
$(GOVULNCHECK): $(LOCALBIN)
	$(call go-install-tool,$(GOVULNCHECK),golang.org/x/vuln/cmd/govulncheck,$(GOVULNCHECK_VERSION))

.PHONY: toolkit-tools-gen
toolkit-tools-gen: $(TOOLKIT_TOOLS_GEN) ## Download toolkit-tools-gen locally if necessary.
$(TOOLKIT_TOOLS_GEN): $(LOCALBIN)
	$(call go-install-tool,$(TOOLKIT_TOOLS_GEN),github.com/emilkje/go-openai-toolkit/cmd/toolkit-tools-gen,$(TOOLKIT_TOOLS_GEN_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef
