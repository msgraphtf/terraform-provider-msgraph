package objectplanmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UseStateForUnconfigured() planmodifier.Object {
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

// PlanModifyObject implements the plan modification logic.
func (m useStateForUnconfiguredModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {

	if !req.ConfigValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}
