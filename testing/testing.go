package main

import (
	"fmt"
    "context"

    azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

func main() {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	
	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
	}

	client , err  := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}


	qparams := users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: []string{"PasswordProfile"},
		},
	}

	result, err := client.Users().ByUserId("0acc8010-50ea-4a54-bd71-ec485d425a74").Get(context.Background(), &qparams)
	if err != nil {
		fmt.Printf("Error getting the user: %v\n", err)
		printOdataError(err)
	}
	fmt.Printf("Found User : %+v\n", *result.GetPasswordProfile().GetForceChangePasswordNextSignIn())

}


func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error:", typed.Error())
		if terr := typed.GetErrorEscaped(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
