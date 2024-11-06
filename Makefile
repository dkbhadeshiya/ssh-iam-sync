NAME=ssh-iam-sync

BUILD_OS= linux darwin
BUILD_ARCH= amd64 arm64

all: clean build packagedeb
build:
	for OS in $(BUILD_OS); do \
		for ARCH in $(BUILD_ARCH); do \
			GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=0 go build -o bin/$(NAME)-$$GITHUB_REF_NAME-$$OS-$$ARCH cmd/$(NAME)/*.go ; \
		done ; \
	done

packagedeb:
	mkdir -p build
	for ARCH in $(BUILD_ARCH); do \
		mkdir -p build/$(NAME)-$$GITHUB_REF_NAME-$$ARCH/usr/bin ; \
		mkdir -p build/$(NAME)-$$GITHUB_REF_NAME-$$ARCH/etc/ssh-iam-sync ; \
		cp config.example.yaml build/$(NAME)-$$GITHUB_REF_NAME-$$ARCH/etc/ssh-iam-sync/config.example.yaml ; \
		cp bin/$(NAME)-$$GITHUB_REF_NAME-linux-$$ARCH build/$(NAME)-$$GITHUB_REF_NAME-$$ARCH/usr/bin/$(NAME) ; \
		chmod +x build/$(NAME)-$$GITHUB_REF_NAME-$$ARCH/usr/bin/$(NAME) ; \
	done
clean:
	go clean
	rm -rf bin/*
	rm -rf build
