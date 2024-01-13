package main

import (
	"context"
	"fmt"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

func main() {

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Printf("Error creating credentials: %s\n", err)
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		fmt.Printf("Error creating client: %s\n", err)
		return
	}

	requestBody := models.NewUser()
	myId := "82d80745-e0a8-4254-9b93-964f8f0a7672"
	myName := ""
	requestBody.SetGivenName(&myName)
	_, err = client.Users().ByUserId(myId).Patch(context.Background(), requestBody, nil)
	if err != nil {
		fmt.Printf("Error patching user\n")
		fmt.Printf("%s\n", err.Error())
	}

	userObject, err := serialization.SerializeToJson(requestBody)
	fmt.Printf("%+v\n", string(userObject))
	fmt.Printf("%+v\n", err)

}
