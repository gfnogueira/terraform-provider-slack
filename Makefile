TEST?=$$(go list ./... |grep -v 'vendor')
SWEEP_DIR?=./slack
SWEEP?=eu-west-1
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
VERSION?=1.0.0
OUTPUT_DIR?=~/.terraform.d/plugins/local/slack/$(VERSION)/linux_amd64

default: testacc build

build:
	@echo "==> Building provider binaries for multiple platforms..."
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/terraform-provider-slack_linux_amd64
	GOOS=linux GOARCH=arm64 go build -o dist/terraform-provider-slack_linux_arm64
	GOOS=darwin GOARCH=amd64 go build -o dist/terraform-provider-slack_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o dist/terraform-provider-slack_darwin_arm64
	GOOS=windows GOARCH=amd64 go build -o dist/terraform-provider-slack_windows_amd64.exe

install: build
	@echo "==> Installing provider plugin..."
	mkdir -p $(OUTPUT_DIR)
	mv dist/terraform-provider-slack_linux_amd64 $(OUTPUT_DIR)/terraform-provider-slack

release:
	@echo "==> Creating release with GoReleaser..."
	goreleaser release --rm-dist

download:
	@echo "==> Download go.mod dependencies"
	@go mod download

tools: download
	@echo "==> Installing tools from tools.go"
	@go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install

lint: tools vet docs
	@echo "==> Checking source code against linters..."
	golangci-lint run -v

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout 60m

test: lint
	go test -v $(TEST)

testacc:
	@echo "==> Running acceptance tests..."
	TF_ACC=1 go test -v $(TEST) $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/)

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || 		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)

docs: docs-lint docscheck

docs-lint:
	@echo "==> Checking docs against linters..."
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli docs/ || (echo; echo "Unexpected issues found in docs Markdown files."; exit 1)
	@terrafmt diff ./docs --check --pattern '*.md' --quiet || (echo; echo "Unexpected differences in docs HCL formatting."; exit 1)

docs-lint-fix:
	@echo "==> Applying automatic docs linter fixes..."
	@docker run -v $(PWD):/markdown 06kellyjac/markdownlint-cli --fix docs/
	@terrafmt fmt ./docs --pattern '*.md'

docscheck:
	@tfproviderdocs check 		-allowed-resource-subcategories-file docs/allowed-subcategories.txt 		-require-resource-subcategory

.PHONY: build download tools lint sweep test testacc vet depscheck docs docs-lint docs-lint-fix docscheck release