qparams := users.UserItemRequestBuilderGetRequestConfiguration{
	QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
		Select: UserProperties[:],
	},
}

var result models.Userable
var err error
if !state.Id.IsNull() {
	result, err = d.client.Users().ByUserId(state.Id.ValueString()).Get(context.Background(), &qparams)
} else if !state.UserPrincipalName.IsNull() {
	result, err = d.client.Users().ByUserId(state.UserPrincipalName.ValueString()).Get(context.Background(), &qparams)
} else {
	resp.Diagnostics.AddError(
		"Missing argument",
		"Either `id` or `user_principal_name` must be supplied.",
	)
	return
}
if err != nil {
	resp.Diagnostics.AddError(
		"Error getting user",
		err.Error(),
	)
	return
}
