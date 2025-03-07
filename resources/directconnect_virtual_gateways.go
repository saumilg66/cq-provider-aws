package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/cloudquery/cq-provider-aws/client"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func DirectconnectVirtualGateways() *schema.Table {
	return &schema.Table{
		Name:         "aws_directconnect_virtual_gateways",
		Resolver:     fetchDirectconnectVirtualGateways,
		Multiplex:    client.AccountRegionMultiplex,
		IgnoreError:  client.IgnoreAccessDeniedServiceDisabled,
		DeleteFilter: client.DeleteAccountRegionFilter,
		Columns: []schema.Column{
			{
				Name:     "account_id",
				Type:     schema.TypeString,
				Resolver: client.ResolveAWSAccount,
			},
			{
				Name:     "region",
				Type:     schema.TypeString,
				Resolver: client.ResolveAWSRegion,
			},
			{
				Name: "virtual_gateway_id",
				Type: schema.TypeString,
			},
			{
				Name: "virtual_gateway_state",
				Type: schema.TypeString,
			},
		},
	}
}

func fetchDirectconnectVirtualGateways(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	var config directconnect.DescribeVirtualGatewaysInput
	c := meta.(*client.Client)
	svc := c.Services().Directconnect
	output, err := svc.DescribeVirtualGateways(ctx, &config, func(options *directconnect.Options) {
		options.Region = c.Region
	})
	if err != nil {
		return err
	}
	res <- output.VirtualGateways
	return nil
}
