default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

gendocs:
	go generate ./...

gends:
	go run ./generate/ $(args); gofmt -w -s -l msgraph/
