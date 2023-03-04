package msgraph

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_enabled": schema.BoolAttribute{
				Required: true,
			},
			"display_name": schema.StringAttribute{
				Required: true,
			},
			"mail_nickname": schema.StringAttribute{
				Required: true,
			},
			"password_profile": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"force_change_password_next_sign_in": schema.BoolAttribute{
						Required: true,
					},
					"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
						Required: true,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
				},
			},
			"user_principal_name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type userResourceModel struct {
	AccountEnabled    types.Bool               `tfsdk:"account_enabled"`
	DisplayName       types.String             `tfsdk:"display_name"`
	MailNickname      types.String             `tfsdk:"mail_nickname"`
	PasswordProfile   *userPasswordProfileModel `tfsdk:"password_profile"`
	UserPrincipalName types.String             `tfsdk:"user_principal_name"`
	Id                types.String             `tfsdk:"id"`
}

type userPasswordProfileModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
	Password                             types.String `tfsdk:"password"`
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Plan
	requestBody := models.NewUser()
	accountEnabled := plan.AccountEnabled.ValueBool()
	requestBody.SetAccountEnabled(&accountEnabled)
	displayName := plan.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)
	mailNickname := plan.MailNickname.ValueString()
	requestBody.SetMailNickname(&mailNickname)

	passwordProfile := models.NewPasswordProfile()
	forceChangePasswordNextSignIn := plan.PasswordProfile.ForceChangePasswordNextSignIn.ValueBool()
	passwordProfile.SetForceChangePasswordNextSignIn(&forceChangePasswordNextSignIn)
	forceChangePasswordNextSignInMfa := plan.PasswordProfile.ForceChangePasswordNextSignInWithMfa.ValueBool()
	passwordProfile.SetForceChangePasswordNextSignInWithMfa(&forceChangePasswordNextSignInMfa)
	password := plan.PasswordProfile.Password.ValueString()
	passwordProfile.SetPassword(&password)
	requestBody.SetPasswordProfile(passwordProfile)

	userPrincipalName := plan.UserPrincipalName.ValueString()
	requestBody.SetUserPrincipalName(&userPrincipalName)

	// Create new User
	result, err := r.client.Users().Post(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			printOdataError(err),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	plan.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: []string{"accountEnabled, displayName, mailNickname, passwordProfile, userPrincipalName, Id"},
		},
	}
	result, err := r.client.UsersById(state.Id.ValueString()).Get(context.Background(), &qparams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting user",
			err.Error(),
		)
	}

	// Overwrite items with refreshed state
	state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	state.DisplayName = types.StringValue(*result.GetDisplayName())
	state.MailNickname = types.StringValue(*result.GetMailNickname())

	passwordProfile := new(userPasswordProfileModel)
	passwordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
	passwordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
	passwordProfile.Password = types.StringNull()
	state.PasswordProfile = passwordProfile

	state.PasswordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
	state.PasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())


	state.UserPrincipalName = types.StringValue(*result.GetUserPrincipalName())
	state.Id = types.StringValue(*result.GetId())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state userResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewUser()
	accountEnabled := plan.AccountEnabled.ValueBool()
	requestBody.SetAccountEnabled(&accountEnabled)
	displayName := plan.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)
	mailNickname := plan.MailNickname.ValueString()
	requestBody.SetMailNickname(&mailNickname)

	passwordProfile := models.NewPasswordProfile()
	forceChangePasswordNextSignIn := plan.PasswordProfile.ForceChangePasswordNextSignIn.ValueBool()
	passwordProfile.SetForceChangePasswordNextSignIn(&forceChangePasswordNextSignIn)
	forceChangePasswordNextSignInMfa := plan.PasswordProfile.ForceChangePasswordNextSignInWithMfa.ValueBool()
	passwordProfile.SetForceChangePasswordNextSignInWithMfa(&forceChangePasswordNextSignInMfa)
	password := plan.PasswordProfile.Password.ValueString()
	passwordProfile.SetPassword(&password)
	requestBody.SetPasswordProfile(passwordProfile)

	userPrincipalName := plan.UserPrincipalName.ValueString()
	requestBody.SetUserPrincipalName(&userPrincipalName)

	_, err := r.client.UsersById(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user: "+plan.Id.ValueString(),
			printOdataError(err),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete user
	err := r.client.UsersById(state.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user: "+state.Id.ValueString(),
			printOdataError(err),
		)
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func printOdataError(err error) string {

	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		if terr := typed.GetError(); terr != nil {
			return fmt.Sprintf("error: %s\ncode: %s\nmsg: %s", typed.Error(), *terr.GetCode(), *terr.GetMessage())
		} else {
			return fmt.Sprintf("error: %s", typed.Error())
		}
	default:
		return fmt.Sprintf("%T > error: %#v", err, err)
	}
}
