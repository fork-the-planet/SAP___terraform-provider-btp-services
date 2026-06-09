// btpservices/provider/cicd/settings/datasource_allowed_spaces.go

package cicdsettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	cicdclient "github.com/SAP/terraform-provider-btp-services/internal/cicd/client"
	"github.com/SAP/terraform-provider-btp-services/internal/shared"
)

var _ datasource.DataSource = &allowedSpacesDataSource{}
var _ datasource.DataSourceWithConfigure = &allowedSpacesDataSource{}

func NewAllowedSpacesDataSource() datasource.DataSource {
	return &allowedSpacesDataSource{}
}

type allowedSpacesDataSource struct {
	cli *cicdclient.CicdClientFacade
}

func (d *allowedSpacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_cicd_allowed_spaces", req.ProviderTypeName)
}

func (d *allowedSpacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the list of allowed spaces in the SAP BTP CI/CD service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the data source.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of allowed Cloud Foundry spaces.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"space_guid": schema.StringAttribute{
							MarkdownDescription: "GUID of the Cloud Foundry space.",
							Computed:            true,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "Human-readable note about why this space is allowed.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *allowedSpacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clients, ok := req.ProviderData.(*shared.ProviderClients)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *shared.ProviderClients, got: %T", req.ProviderData),
		)
		return
	}
	if clients.Cicd == nil {
		resp.Diagnostics.AddError(
			"Missing CI/CD Configuration",
			"A cicd{} block must be configured in the provider to use CI/CD data sources.",
		)
		return
	}
	d.cli = clients.Cicd
}

func (d *allowedSpacesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	result, err := d.cli.AllowedSpaces.Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error Listing Allowed Spaces", err.Error())
		return
	}

	entries := make([]allowedSpaceEntry, 0, len(result.AllowedSpaces))
	for _, s := range result.AllowedSpaces {
		entries = append(entries, allowedSpaceEntry{
			SpaceGUID: types.StringValue(s.SpaceGUID),
			Comment:   types.StringValue(s.Comment),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, allowedSpacesDSModel{
		ID:     types.StringValue("allowed_spaces"),
		Values: entries,
	})...)
}

type allowedSpacesDSModel struct {
	ID     types.String        `tfsdk:"id"`
	Values []allowedSpaceEntry `tfsdk:"values"`
}
