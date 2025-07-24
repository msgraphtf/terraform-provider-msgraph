build:
	git clean -f msgraph/; go run ./generate/ $(args); gofmt -w -s -l msgraph/; go build

install: build
	go install

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: testacc generate
