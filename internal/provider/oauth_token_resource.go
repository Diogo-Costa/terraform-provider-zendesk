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
	_ resource.Resource                = &OAuthTokenResource{}
	_ resource.ResourceWithImportState = &OAuthTokenResource{}
)

func NewOAuthTokenResource() resource.Resource {
	return &OAuthTokenResource{}
}

type OAuthTokenResource struct {
	client *Client
}

type OAuthTokenResourceModel struct {
	ID        types.String   `tfsdk:"id"`
	ClientID  types.String   `tfsdk:"client_id"`
	Scopes    []types.String `tfsdk:"scopes"`
	FullToken types.String   `tfsdk:"full_token"`
	ExpiresAt types.String   `tfsdk:"expires_at"`
}

func (r *OAuthTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_token"
}

func (r *OAuthTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Zendesk OAuth token.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the OAuth token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The ID of the OAuth client.",
				Required:    true,
			},
			"scopes": schema.ListAttribute{
				Description: "The scopes granted to the OAuth token.",
				Required:    true,
				ElementType: types.StringType,
			},
			"full_token": schema.StringAttribute{
				Description: "The full OAuth token value (only available after creation).",
				Computed:    true,
				Sensitive:   true,
			},
			"expires_at": schema.StringAttribute{
				Description: "The expiration date of the token in ISO 8601 format (e.g., '2024-12-31T23:59:59Z'). If not set, the token will not expire.",
				Optional:    true,
			},
		},
	}
}

func (r *OAuthTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OAuthTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OAuthTokenResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientID, err := strconv.ParseInt(plan.ClientID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing Client ID",
			fmt.Sprintf("Could not parse client ID: %v", err),
		)
		return
	}

	scopes := make([]string, 0, len(plan.Scopes))
	for _, scope := range plan.Scopes {
		scopes = append(scopes, scope.ValueString())
	}

	token, err := r.client.CreateOAuthToken(clientID, scopes, plan.ExpiresAt.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating OAuth Token",
			fmt.Sprintf("Could not create OAuth token: %v", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(token.ID, 10))
	plan.FullToken = types.StringValue(token.FullToken)
	plan.ExpiresAt = types.StringValue(token.ExpiresAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *OAuthTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OAuthTokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing OAuth Token ID",
			fmt.Sprintf("Could not parse OAuth token ID: %v", err),
		)
		return
	}

	token, err := r.client.ReadOAuthToken(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading OAuth Token",
			fmt.Sprintf("Could not read OAuth token: %v", err),
		)
		return
	}

	if token == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ClientID = types.StringValue(strconv.FormatInt(token.ClientID, 10))
	state.Scopes = make([]types.String, 0, len(token.Scopes))
	for _, scope := range token.Scopes {
		state.Scopes = append(state.Scopes, types.StringValue(scope))
	}
	state.ExpiresAt = types.StringValue(token.ExpiresAt)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *OAuthTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"The Zendesk API does not support updating OAuth tokens. To change the configuration, you must create a new token.",
	)
}

func (r *OAuthTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OAuthTokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Parsing OAuth Token ID",
			fmt.Sprintf("Could not parse OAuth token ID: %v", err),
		)
		return
	}

	err = r.client.DeleteOAuthToken(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting OAuth Token",
			fmt.Sprintf("Could not delete OAuth token: %v", err),
		)
		return
	}
}

func (r *OAuthTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
} 