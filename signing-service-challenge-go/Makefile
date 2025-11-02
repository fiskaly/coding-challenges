GOCACHE ?= $(CURDIR)/.cache
GOMODCACHE ?= $(CURDIR)/.gomodcache
GOENV := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)

.PHONY: run test clean-cache

run: $(GOCACHE) $(GOMODCACHE)
	@$(GOENV) go run ./main.go

test: $(GOCACHE) $(GOMODCACHE)
	@$(GOENV) go test ./...

clean-cache:
	@rm -rf $(GOCACHE) $(GOMODCACHE)

$(GOCACHE) $(GOMODCACHE):
	@mkdir -p $@
