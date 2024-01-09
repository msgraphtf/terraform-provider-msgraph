package listplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UseStateForUnconfigured() planmodifier.List {
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

// PlanModifyList implements the plan modification logic.
func (m useStateForUnconfiguredModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {

	if !req.ConfigValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}
