// btpservices/provider/cicd/settings/list_resource_allowed_spaces.go

package cicdsettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	cicdclient "github.com/SAP/terraform-provider-btp-services/internal/cicd/client"
	"github.com/SAP/terraform-provider-btp-services/internal/shared"
)

var _ list.ListResourceWithConfigure = &allowedSpacesListResource{}

type allowedSpacesListResource struct {
	cli *cicdclient.CicdClientFacade
}

func NewAllowedSpacesListResource() list.ListResource {
	return &allowedSpacesListResource{}
}

func (r *allowedSpacesListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_cicd_allowed_spaces", req.ProviderTypeName)
}

func (r *allowedSpacesListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			"A cicd{} block must be configured in the provider to use CI/CD resources.",
		)
		return
	}
	r.cli = clients.Cicd
}

func (r *allowedSpacesListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "This list resource discovers all allowed spaces in the SAP BTP CI/CD service.",
	}
}

func (r *allowedSpacesListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	apiResult, err := r.cli.AllowedSpaces.Get(ctx)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Listing Allowed Spaces",
			fmt.Sprintf("Failed to list allowed spaces: %s", err),
		)
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		item := req.NewListResult(ctx)
		item.Identity.SetAttribute(ctx, path.Root("id"), "allowed_spaces")

		if req.IncludeResource {
			model := allowedSpacesValueFrom(*apiResult)
			item.Diagnostics.Append(item.Resource.Set(ctx, model)...)
		}

		push(item)
	}
}
