tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/katbyte/terrafmt@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.41.1

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

goimports:
	@echo "==> Fixing source code with gofmt..."
	# This logic should match the search logic in scripts/gofmtcheck.sh
	find . -name '*.go' | grep -v vendor | xargs goimports -w

lint:
	golangci-lint run ./...

.PHONY: tools fmt goimports lint