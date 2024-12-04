// Copyright (c) HashiCorp, Inc.

// Description: This file contains the implementation of the FreeipaHostDataSource data source.

package provider

import (
	"context"
	"fmt"
	"github.com/ccin2p3/go-freeipa/freeipa"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FreeipaHostDataSource{}
var _ datasource.DataSourceWithConfigure = &FreeipaHostDataSource{}

func NewFreeipaHostDataSource() datasource.DataSource {
	return &FreeipaHostDataSource{}
}

type FreeipaHostDataSource struct {
	client *freeipa.Client
}

type FreeipaHostDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	Fqdn     types.String `tfsdk:"fqdn"`
	Hostname types.String `tfsdk:"hostname"`
}

func (d *FreeipaHostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (d *FreeipaHostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Freeipa host data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the host",
				Computed:            true,
			},
			"fqdn": schema.StringAttribute{
				MarkdownDescription: "Fqdn of the host",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the host",
				Computed:            true,
			},
		},
	}
}

func (d *FreeipaHostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*freeipa.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *freeipa.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *FreeipaHostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FreeipaHostDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := d.client.HostShow(
		&freeipa.HostShowArgs{
			Fqdn: data.Fqdn.ValueString(),
		}, &freeipa.HostShowOptionalArgs{})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", data.Fqdn.String(), err))
		return
	}

	data.Id = types.StringValue(host.Result.Fqdn)
	data.Fqdn = types.StringValue(host.Result.Fqdn)
	data.Hostname = types.StringValue(host.Result.Fqdn)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
