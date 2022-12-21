.PHONY: clean
clean:
	@rm -rf dist

.PHONY: setup
setup:
	@./scripts/setup.sh

.PHONY: build
build: $(CONTEXT_FILES)
	@./scripts/build.sh
