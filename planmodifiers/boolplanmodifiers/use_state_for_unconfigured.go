package boolplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UseStateForUnconfigured() planmodifier.Bool {
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

// PlanModifyBool implements the plan modification logic.
func (m useStateForUnconfiguredModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {

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
