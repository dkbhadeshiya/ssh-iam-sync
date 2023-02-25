NAME=ssh-iam-sync

BUILD_OS= linux darwin
BUILD_ARCH= amd64 arm64

all: clean build
build:
	for OS in $(BUILD_OS); do \
		for ARCH in $(BUILD_ARCH); do \
			GOOS=$$OS GOARCH=$$ARCH go build -o bin/$(NAME)-$$GITHUB_REF_NAME-$$OS-$$ARCH cmd/$(NAME)/*.go ; \
		done ; \
	done

clean:
	go clean
	rm -rf bin/*
