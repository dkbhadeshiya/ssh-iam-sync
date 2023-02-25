all:  vet lint build

vet:
    go vet ./...

lint:
    go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

build:
    go build -o bin/ssh-iam-sync ./cmd/ssh-iam-sync
