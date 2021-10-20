tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/katbyte/terrafmt@latest

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

goimports:
	@echo "==> Fixing source code with gofmt..."
	# This logic should match the search logic in scripts/gofmtcheck.sh
	find . -name '*.go' | grep -v vendor | xargs goimports -w

.PHONY: tools fmt goimports