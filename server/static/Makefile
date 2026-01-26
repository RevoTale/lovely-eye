.PHONY: build typecheck

GO ?= go

TSC ?= bunx tsc

build:
	$(GO) run build.go

typecheck:
	$(TSC) -p tsconfig.json
