package stringplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UseStateForUnconfigured() planmodifier.String {
	return useStateForUnconfiguredModifier{}
}

// useStateForUnknownModifier implements the plan modifier.
type useStateForUnconfiguredModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m useStateForUnconfiguredModifier) Description(_ context.Context) string {
	return "If unconfigured, this attribute will use the value in state."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useStateForUnconfiguredModifier) MarkdownDescription(_ context.Context) string {
	return "If unconfigured, this attribute will use the value in state."
}

// PlanModifyString implements the plan modification logic.
func (m useStateForUnconfiguredModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {

	// Do nothing if resource is being created
	if req.State.Raw.IsNull() {
		return
	}

	// Do nothing if configuration is not null
	if !req.ConfigValue.IsNull() {
		return
	}

	// If resource is being updated, and config is null, use state value
	resp.PlanValue = req.StateValue
}
