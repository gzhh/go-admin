.PHONY: init
init:
	@echo "This is init step"

.PHONY: lint
lint:
	./scripts/lint.sh

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./bin/adminapp -a -v -ldflags '-s -w' -tags=jsoniter cmd/admin/main.go
