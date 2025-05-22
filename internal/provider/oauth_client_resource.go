package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &OAuthClientResource{}
	_ resource.ResourceWithImportState = &OAuthClientResource{}
)

func NewOAuthClientResource() resource.Resource {
	return &OAuthClientResource{}
}

type OAuthClientResource struct {
	client *Client
}

type OAuthClientResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Identifier  types.String `tfsdk:"identifier"`
	Kind        types.String `tfsdk:"kind"`
	Description types.String `tfsdk:"description"`
}

func (r *OAuthClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_client"
}

func (r *OAuthClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Zendesk OAuth client.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the OAuth client.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the OAuth client.",
				Required:    true,
			},
			"identifier": schema.StringAttribute{
				Description: "The unique identifier of the OAuth client.",
				Required:    true,
			},
			"kind": schema.StringAttribute{
				Description: "The kind of OAuth client (e.g., 'public').",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the OAuth client.",
				Optional:    true,
			},
		},
	}
}

func (r *OAuthClientResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *OAuthClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OAuthClientResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.client.CreateOAuthClient(
		plan.Name.ValueString(),
		plan.Identifier.ValueString(),
		plan.Kind.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating OAuth Client",
			fmt.Sprintf("Could not create OAuth client: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(client.ID, 10))
	plan.Description = types.StringValue(client.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *OAuthClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OAuthClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing OAuth Client ID",
			fmt.Sprintf("Could not parse OAuth client ID: %v", err),
		)
		return
	}

	client, err := r.client.ReadOAuthClient(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading OAuth Client",
			fmt.Sprintf("Could not read OAuth client: %v", err),
		)
		return
	}

	if client == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(client.Name)
	state.Identifier = types.StringValue(client.Identifier)
	state.Kind = types.StringValue(client.Kind)
	state.Description = types.StringValue(client.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *OAuthClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"The Zendesk API does not support updating OAuth clients. To change the configuration, you must create a new client.",
	)
}

func (r *OAuthClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OAuthClientResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing OAuth Client ID",
			fmt.Sprintf("Could not parse OAuth client ID: %v", err),
		)
		return
	}

	err = r.client.DeleteOAuthClient(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting OAuth Client",
			fmt.Sprintf("Could not delete OAuth client: %v", err),
		)
		return
	}
}

func (r *OAuthClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
} 