package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"terraform-provider-msgraph/generate/openapi"
	"terraform-provider-msgraph/generate/transform"
)

func getAugment(pathname string) transform.TemplateAugment {
	pathObject := openapi.GetPath(pathname)

	pathFields := strings.Split(pathObject.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName := strings.ToLower(pathFields[0])

	// Open augment file if available
	var err error = nil
	augment := transform.TemplateAugment{}
	augmentFile, err := os.ReadFile("generate/augment/" + packageName + "/" + getBlockName(pathname) + ".yaml")
	if err == nil {
		yaml.Unmarshal(augmentFile, &augment)
	}

	return augment

}

func getBlockName(pathname string) string {

	pathObject := openapi.GetPath(pathname)
	pathFields := strings.Split(pathObject.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array

	// Generate name of the terraform block
	blockName := ""
	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := transform.PathFieldName(p)
				blockName += pLeft
			} else {
				blockName += p
			}
		}
	} else {
		blockName = pathFields[0]
	}

	return blockName
}

func main() {

	if len(os.Args) > 1 {
		pathObject := openapi.GetPath(os.Args[1])
		blockName := getBlockName(os.Args[1])
		augment := getAugment(os.Args[1])
		generateDataSource(pathObject, blockName, augment)
		generateModel(pathObject, blockName, augment)
		if pathObject.Patch.Summary != "" {
			generateResource(pathObject, blockName, augment)
		}
	} else {

		// TODO: Change from using paths to using tags and/or operation IDs.
		// This should help to remove duplicate paths, and duplicate model stuff

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
			pathObject := openapi.GetPath(path)
			blockName := getBlockName(path)
			augment := getAugment(path)
			generateDataSource(pathObject, blockName, augment)
			generateModel(pathObject, blockName, augment)
			if pathObject.Patch.Summary != "" && pathObject.Delete.Summary != "" {
				generateResource(pathObject, blockName, augment)
			}
		}

	}

}
