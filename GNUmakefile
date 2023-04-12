default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

gendocs:
	go generate ./main.go

gends:
	go run ./template/template.go $(args); go fmt ./template/out/user_data_source.go
