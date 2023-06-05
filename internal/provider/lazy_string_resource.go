package provider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
  _ resource.Resource = &lazyStringResource{};
  _ resource.ResourceWithModifyPlan = &lazyStringResource{};
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
        Computed: true,
      },
      "last_updated": schema.StringAttribute{
        Computed: true,
      },
      "result": schema.StringAttribute{
        Computed: true,
      },
      "initially": schema.StringAttribute{
        Required: true,
      },
      "explicitly": schema.StringAttribute{
        Optional: true,
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
  var state_id types.String
  var state_result types.String
  var state_last_updated types.String
  resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
  resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &state_id)...)
  resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("result"), &state_result)...)
  resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("last_updated"), &state_last_updated)...)

  if resp.Diagnostics.HasError() {
    return
  }

  if state_id != types.StringNull() {
    plan.ID = state_id
  }

  if plan.Explicitly != types.StringNull() {
    plan.Result = plan.Explicitly
  } else if state_result != types.StringNull() {
    plan.Result = state_result
    plan.LastUpdated = state_last_updated
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

func (r *lazyStringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
  }
  
  resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
  if resp.Diagnostics.HasError() {
    return
  }
}

func (r *lazyStringResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

