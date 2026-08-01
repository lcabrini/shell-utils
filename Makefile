APPS := ankify

all: build

.PHONY: build
build: pre-build
	for app in $(APPS); do \
		go build ./cmd/$$app; \
	done

.PHONY: pre-build
pre-build:
	go mod tidy

.PHONY: install
insall:
	for app in $(APP); do \
		go install ./cmd/$$app; \
	done

.PHONY: clean
clean:
	go clean
	rm -f $(APPS)
