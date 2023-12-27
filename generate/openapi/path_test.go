package openapi

import (
	"fmt"
	"testing"
)

func TestGetPath(t *testing.T) {

	path := GetPath("/users")
	fmt.Printf("PATH:           %s\n", path.Path)
	fmt.Printf("DESCRIPTION:    %s\n", path.Description)
	fmt.Printf("PARAMETERS:     %v\n", path.Parameters)
	fmt.Printf("GET RETURNS:    %+v\n", path.Get)
	fmt.Printf("POST RETURNS:   %+v\n", path.Post)
	fmt.Printf("PATCH RETURNS:  %+v\n", path.Patch)
	fmt.Printf("DELETE RETURNS: %+v\n", path.Delete)

}
