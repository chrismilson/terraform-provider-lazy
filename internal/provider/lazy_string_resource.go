package provider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &lazyStringResource{}
	_ resource.ResourceWithModifyPlan  = &lazyStringResource{}
	_ resource.ResourceWithImportState = &lazyStringResource{}
)

type lazyStringResource struct{}

func NewLazyStringResource() resource.Resource {
	return &lazyStringResource{}
}

func (r *lazyStringResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_string"
}

func (r *lazyStringResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: `This is set to a random value at create time.`,
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: `This is updated to the current time whenever ` + "`result`" + ` is changing.`,
				Computed:    true,
			},
			"result": schema.StringAttribute{
				Description: `This is the canonical value for the resource. Depends on whether a value is (or has been in the past) supplied explicitly, or initially.`,
				Computed:    true,
			},
			"initially": schema.StringAttribute{
				Description: `This will determine the value of the ` + "`result`" + ` on create if no value is supplied explicitly. In general this should be a static value to avoid unnecessary terraform changes.`,
				Computed:    true,
				Optional:    true,
			},
			"explicitly": schema.StringAttribute{
				Description: `When set, this will determine the value of ` + "`result`" + `. If unset, it will be maintained at its previous value to avoid changes in the plan.`,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

type lazyStringResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Result      types.String `tfsdk:"result"`
	Initially   types.String `tfsdk:"initially"`
	Explicitly  types.String `tfsdk:"explicitly"`
}

func (r *lazyStringResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan lazyStringResourceModel
	var state lazyStringResourceModel

	if req.Plan.Raw.IsNull() {
		// About to delete
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if req.State.Raw.IsNull() {
		// About to create
		state = lazyStringResourceModel{}
	} else {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID != types.StringNull() {
		// ID should never change
		plan.ID = state.ID
	}

	// Calculate "initially"
	if plan.Initially == types.StringNull() || plan.Initially == types.StringUnknown() {
		// When omitted, there should be no changes in the plan
		plan.Initially = state.Initially
	}

	// Calculate "explicitly"
	if plan.Explicitly == types.StringNull() || plan.Explicitly == types.StringUnknown() {
		// When omitted, there should be no changes in the plan
		plan.Explicitly = state.Explicitly
	}

	// Calculate "result"
	if plan.Explicitly != types.StringNull() {
		plan.Result = plan.Explicitly
	} else {
		plan.Result = plan.Initially
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *lazyStringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan lazyStringResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d", rand.Int()))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *lazyStringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lazyStringResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Result != state.Result {
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	} else {
		plan.LastUpdated = state.LastUpdated
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *lazyStringResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	plan := lazyStringResourceModel{
		ID:          types.StringValue(fmt.Sprintf("%d", rand.Int())),
		LastUpdated: types.StringValue(time.Now().Format(time.RFC850)),
		Explicitly:  types.StringValue(req.ID),
		Initially:   types.StringValue(req.ID),
		Result:      types.StringValue(req.ID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *lazyStringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *lazyStringResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
