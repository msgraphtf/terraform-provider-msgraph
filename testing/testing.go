package main

import (
	"context"
	"fmt"
	"reflect"

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

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
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

	client.Users().ByUserId("").Activities().Get()

	fmt.Printf("%+v\n", reflect.ValueOf(client.Users().ByUserId("0acc8010-50ea-4a54-bd71-ec485d425a74")).MethodByName("Get").Type())

	value, _ := reflect.TypeOf(client.Users()).MethodByName("ByUserId")
	value, _ = value.Type.Out(0).MethodByName("Get")
	method := value.Type.In(2).Elem()
	value2, _ := method.FieldByName("QueryParameters")
	fmt.Printf("%+v\n", value2.Type)

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
