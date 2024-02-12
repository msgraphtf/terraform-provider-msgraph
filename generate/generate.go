package main

import (
	"os"
	"strings"
	"gopkg.in/yaml.v3"

	"terraform-provider-msgraph/generate/openapi"
)

var pathObject openapi.OpenAPIPathObject
var blockName string
var augment templateAugment
var input templateInput

func setGlobals(pathname string) {
	pathObject = openapi.GetPath(pathname)

	pathFields := strings.Split(pathObject.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName := strings.ToLower(pathFields[0])

	// Generate name of the terraform block
	blockName = ""
	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := pathFieldName(p)
				blockName += pLeft
			} else {
				blockName += p
			}
		}
	} else {
		blockName = pathFields[0]
	}

	// Open augment file if available
	var err error = nil
	augment = templateAugment{}
	augmentFile, err := os.ReadFile("generate/augment/" + packageName + "/" + blockName + ".yaml")
	if err == nil {
		yaml.Unmarshal(augmentFile, &augment)
	}

}

func main() {

	if len(os.Args) > 1 {
		setGlobals(os.Args[1])
		generateDataSource()
		generateModel()
	} else {

		knownGoodPaths := [...]string{
			"/applications",
			"/applications/{application-id}",
			"/devices",
			"/devices/{device-id}",
			"/groups",
			"/groups/{group-id}",
			"/servicePrincipals",
			"/servicePrincipals/{servicePrincipal-id}",
			"/sites",
			"/sites/{site-id}",
			"/teams/{team-id}",
			"/users",
			"/users/{user-id}",
		}

		for _, path := range knownGoodPaths {
			setGlobals(path)
			generateDataSource()
			generateModel()
		}

	}

}
