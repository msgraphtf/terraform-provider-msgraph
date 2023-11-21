package openapi

import (
	"testing"
)

func TestGetPath(t *testing.T) {

	GetPath("/users/{user-id}", "../../msgraph-metadata/openapi/v1.0/openapi.yaml")

}
