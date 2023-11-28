package openapi

import (
	"testing"
	"fmt"
)

func TestGetPath(t *testing.T) {

	path := GetPath("/users/{user-id}")
	fmt.Printf("PATH:        %s\n", path.Path)
	fmt.Printf("DESCRIPTION: %s\n", path.Description)
	fmt.Printf("PARAMETERS:  %v\n", path.Parameters)
	fmt.Printf("GET RETURNS: %s\n", path.Get.Response.Title)

}
