package provider

import (
	"context"
	"fmt"
	"github.com/ccin2p3/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"terraform-provider-freeipa/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FreeipaHostResource{}
var _ resource.ResourceWithImportState = &FreeipaHostResource{}

func NewFreeipaHostResource() resource.Resource {
	return &FreeipaHostResource{}
}

type FreeipaHostResource struct {
	client *freeipa.Client
}

type FreeipaHostResourceModel struct {
	Fqdn        types.String `tfsdk:"fqdn"`
	Description types.String `tfsdk:"description"`
	Force       types.Bool   `tfsdk:"force"`
	NoReverse   types.Bool   `tfsdk:"noreverse"`
	Id          types.String `tfsdk:"id"`
}

func (r *FreeipaHostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *FreeipaHostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Freeipa host resource",

		Attributes: map[string]schema.Attribute{
			"fqdn": schema.StringAttribute{
				MarkdownDescription: "Fqdn of the host",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the host",
				Optional:            true,
			},
			"force": schema.BoolAttribute{
				MarkdownDescription: "Force the operation of host creation irrespective of the dns existence",
				Optional:            true,
			},
			"noreverse": schema.BoolAttribute{
				MarkdownDescription: "Do not create reverse DNS record",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "host identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *FreeipaHostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*freeipa.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *FreeipaHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FreeipaHostResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Fqdn.IsUnknown() {
		resp.Diagnostics.AddError("Missing fqdn", "Fqdn is required to create a host")
	}
	if data.Force.IsUnknown() {
		data.Force = types.BoolValue(true) // setting default value as true
	}
	if data.NoReverse.IsUnknown() {
		data.NoReverse = types.BoolValue(true) // setting default value as true
	}
	if data.Description.IsUnknown() {
		data.Description = types.StringValue("") // setting default value as empty string
	}

	host, err := r.client.HostAdd(&freeipa.HostAddArgs{
		Fqdn: data.Fqdn.ValueString(),
	}, &freeipa.HostAddOptionalArgs{
		Description: utils.RefString(data.Description.ValueString()),
		Force:       utils.RefBool(data.Force.ValueBool()),
		NoReverse:   utils.RefBool(data.NoReverse.ValueBool()),
	})
	if err != nil {
		log.Fatal(err)
	}

	data.Id = types.StringValue(host.Result.Fqdn)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("created host: %s", data.Fqdn.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FreeipaHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FreeipaHostResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := r.client.HostShow(
		&freeipa.HostShowArgs{
			Fqdn: state.Fqdn.ValueString(),
		}, &freeipa.HostShowOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", state.Fqdn.String(), err))
		return
	}

	state.Id = types.StringValue(host.Result.Fqdn)
	state.Fqdn = types.StringValue(host.Result.Fqdn)
	state.Description = types.StringValue(*host.Result.Description)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FreeipaHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FreeipaHostResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if state.Description.ValueString() != plan.Description.ValueString() {
		host, err := r.client.HostMod(
			&freeipa.HostModArgs{
				Fqdn: state.Fqdn.ValueString(),
			}, &freeipa.HostModOptionalArgs{
				Description: utils.RefString(plan.Description.ValueString()),
			})
		if err != nil {
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to update %s, got error: %s", plan.Fqdn.String(), err))
			return
		}
		state.Id = types.StringValue(host.Result.Fqdn)
		state.Fqdn = types.StringValue(host.Result.Fqdn)
		state.Description = types.StringValue(*host.Result.Description)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FreeipaHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FreeipaHostResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.HostDel(
		&freeipa.HostDelArgs{
			Fqdn: []string{data.Fqdn.ValueString()},
		}, &freeipa.HostDelOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s\nSkipping Delete operation in ipa-server and continuing state removal!!", data.Fqdn.String(), err)) // skipping delete operation in ipa-server and continuing state removal
		return
	}
}

func (r *FreeipaHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
