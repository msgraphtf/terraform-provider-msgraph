default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

gendocs:
	go generate ./...

build:
	git clean -f msgraph/; go run ./generate/ $(args); gofmt -w -s -l msgraph/; go build

install: build
	go install
