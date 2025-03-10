package resources

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	elbv1 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/cloudquery/cq-provider-aws/client"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func Elbv1LoadBalancers() *schema.Table {
	return &schema.Table{
		Name:         "aws_elbv1_load_balancers",
		Resolver:     fetchElbv1LoadBalancers,
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
				Name:     "attributes_access_log_enabled",
				Type:     schema.TypeBool,
				Resolver: resolveElbv1loadBalancerAttributesAccessLogEnabled,
			},
			{
				Name:     "attributes_access_log_s3_bucket_name",
				Type:     schema.TypeString,
				Resolver: resolveElbv1loadBalancerAttributesAccessLogS3BucketName,
			},
			{
				Name:     "attributes_access_log_s3_bucket_prefix",
				Type:     schema.TypeString,
				Resolver: resolveElbv1loadBalancerAttributesAccessLogS3BucketPrefix,
			},
			{
				Name:     "attributes_access_log_emit_interval",
				Type:     schema.TypeInt,
				Resolver: resolveElbv1loadBalancerAttributesAccessLogEmitInterval,
			},
			{
				Name:     "attributes_connection_settings_idle_timeout",
				Type:     schema.TypeInt,
				Resolver: resolveElbv1loadBalancerAttributesConnectionSettingsIdleTimeout,
			},
			{
				Name:     "attributes_cross_zone_load_balancing_enabled",
				Type:     schema.TypeBool,
				Resolver: resolveElbv1loadBalancerAttributesCrossZoneLoadBalancingEnabled,
			},
			{
				Name:     "attributes_connection_draining_enabled",
				Type:     schema.TypeBool,
				Resolver: resolveElbv1loadBalancerAttributesConnectionDrainingEnabled,
			},
			{
				Name:     "attributes_connection_draining_timeout",
				Type:     schema.TypeInt,
				Resolver: resolveElbv1loadBalancerAttributesConnectionDrainingTimeout,
			},
			{
				Name:     "attributes_additional_attributes",
				Type:     schema.TypeJSON,
				Resolver: resolveElbv1loadBalancerAttributesAdditionalAttributes,
			},
			{
				Name: "tags",
				Type: schema.TypeJSON,
			},
			{
				Name: "availability_zones",
				Type: schema.TypeStringArray,
			},
			{
				Name: "canonical_hosted_zone_name",
				Type: schema.TypeString,
			},
			{
				Name:     "canonical_hosted_zone_name_id",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("CanonicalHostedZoneNameID"),
			},
			{
				Name: "created_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name:     "dns_name",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("DNSName"),
			},
			{
				Name:     "health_check_healthy_threshold",
				Type:     schema.TypeInt,
				Resolver: schema.PathResolver("HealthCheck.HealthyThreshold"),
			},
			{
				Name:     "health_check_interval",
				Type:     schema.TypeInt,
				Resolver: schema.PathResolver("HealthCheck.Interval"),
			},
			{
				Name:     "health_check_target",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("HealthCheck.Target"),
			},
			{
				Name:     "health_check_timeout",
				Type:     schema.TypeInt,
				Resolver: schema.PathResolver("HealthCheck.Timeout"),
			},
			{
				Name:     "health_check_unhealthy_threshold",
				Type:     schema.TypeInt,
				Resolver: schema.PathResolver("HealthCheck.UnhealthyThreshold"),
			},
			{
				Name:     "instances",
				Type:     schema.TypeStringArray,
				Resolver: resolveElbv1loadBalancerInstances,
			},
			{
				Name: "load_balancer_name",
				Type: schema.TypeString,
			},
			{
				Name:     "other_policies",
				Type:     schema.TypeStringArray,
				Resolver: schema.PathResolver("Policies.OtherPolicies"),
			},
			{
				Name: "scheme",
				Type: schema.TypeString,
			},
			{
				Name: "security_groups",
				Type: schema.TypeStringArray,
			},
			{
				Name:     "source_security_group_name",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("SourceSecurityGroup.GroupName"),
			},
			{
				Name:     "source_security_group_owner_alias",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("SourceSecurityGroup.OwnerAlias"),
			},
			{
				Name: "subnets",
				Type: schema.TypeStringArray,
			},
			{
				Name:     "vpc_id",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("VPCId"),
			},
		},
		Relations: []*schema.Table{
			{
				Name:     "aws_elbv1_load_balancer_backend_server_descriptions",
				Resolver: fetchElbv1LoadBalancerBackendServerDescriptions,
				Columns: []schema.Column{
					{
						Name:     "load_balancer_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name: "instance_port",
						Type: schema.TypeInt,
					},
					{
						Name: "policy_names",
						Type: schema.TypeStringArray,
					},
				},
			},
			{
				Name:     "aws_elbv1_load_balancer_listeners",
				Resolver: fetchElbv1LoadBalancerListeners,
				Columns: []schema.Column{
					{
						Name:     "load_balancer_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name:     "listener_instance_port",
						Type:     schema.TypeInt,
						Resolver: schema.PathResolver("Listener.InstancePort"),
					},
					{
						Name:     "listener_load_balancer_port",
						Type:     schema.TypeInt,
						Resolver: schema.PathResolver("Listener.LoadBalancerPort"),
					},
					{
						Name:     "listener_protocol",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Listener.Protocol"),
					},
					{
						Name:     "listener_instance_protocol",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Listener.InstanceProtocol"),
					},
					{
						Name:     "listener_ssl_certificate_id",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Listener.SSLCertificateId"),
					},
					{
						Name: "policy_names",
						Type: schema.TypeStringArray,
					},
				},
			},
			{
				Name:     "aws_elbv1_load_balancer_policies_app_cookie_stickiness_policies",
				Resolver: fetchElbv1LoadBalancerPoliciesAppCookieStickinessPolicies,
				Columns: []schema.Column{
					{
						Name:     "load_balancer_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name: "cookie_name",
						Type: schema.TypeString,
					},
					{
						Name: "policy_name",
						Type: schema.TypeString,
					},
				},
			},
			{
				Name:     "aws_elbv1_load_balancer_policies_lb_cookie_stickiness_policies",
				Resolver: fetchElbv1LoadBalancerPoliciesLbCookieStickinessPolicies,
				Columns: []schema.Column{
					{
						Name:     "load_balancer_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name: "cookie_expiration_period",
						Type: schema.TypeBigInt,
					},
					{
						Name: "policy_name",
						Type: schema.TypeString,
					},
				},
			},
			{
				Name:     "aws_elbv1_load_balancer_policies",
				Resolver: fetchElbv1LoadBalancerPolicies,
				Columns: []schema.Column{
					{
						Name:     "load_balancer_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name:     "policy_attribute_descriptions",
						Type:     schema.TypeJSON,
						Resolver: resolveElbv1loadBalancerPolicyPolicyAttributeDescriptions,
					},
					{
						Name: "policy_name",
						Type: schema.TypeString,
					},
					{
						Name: "policy_type_name",
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================
func fetchElbv1LoadBalancers(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {

	c := meta.(*client.Client)
	svc := c.Services().ELBv1
	processLoadBalancers := func(loadBalancers []types.LoadBalancerDescription) error {
		tagsCfg := &elbv1.DescribeTagsInput{LoadBalancerNames: make([]string, 0, len(loadBalancers))}
		for _, lb := range loadBalancers {
			tagsCfg.LoadBalancerNames = append(tagsCfg.LoadBalancerNames, *lb.LoadBalancerName)
		}
		tagsResponse, err := svc.DescribeTags(ctx, tagsCfg)
		if err != nil {
			return err
		}
		for _, lb := range loadBalancers {

			tags := getTagsByLoadBalancerName(*lb.LoadBalancerName, tagsResponse.TagDescriptions)

			loadBalancerAttributes, err := svc.DescribeLoadBalancerAttributes(ctx, &elbv1.DescribeLoadBalancerAttributesInput{LoadBalancerName: lb.LoadBalancerName})
			if err != nil {
				return err
			}

			wrapper := ELBv1LoadBalancerWrapper{
				LoadBalancerDescription: lb,
				Tags:                    make(map[string]interface{}, len(tags)),
				Attributes:              loadBalancerAttributes.LoadBalancerAttributes,
			}

			for _, t := range tags {
				wrapper.Tags[*t.Key] = t.Value
			}
			res <- wrapper
		}
		return nil
	}

	var config elbv1.DescribeLoadBalancersInput
	for {
		response, err := svc.DescribeLoadBalancers(ctx, &config, func(options *elbv1.Options) {
			options.Region = c.Region
		})
		if err != nil {
			return err
		}

		for i := 0; i < len(response.LoadBalancerDescriptions); i += 20 {
			end := i + 20

			if end > len(response.LoadBalancerDescriptions) {
				end = len(response.LoadBalancerDescriptions)
			}
			loadBalancers := response.LoadBalancerDescriptions[i:end]
			err := processLoadBalancers(loadBalancers)
			if err != nil {
				return err
			}
		}

		if aws.ToString(response.NextMarker) == "" {
			break
		}
		config.Marker = response.NextMarker
	}

	return nil
}

func resolveElbv1loadBalancerAttributesAccessLogEnabled(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.AccessLog == nil {
		return nil
	}
	return resource.Set(c.Name, r.Attributes.AccessLog.Enabled)
}
func resolveElbv1loadBalancerAttributesAccessLogS3BucketName(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.AccessLog == nil {
		return nil
	}
	return resource.Set(c.Name, r.Attributes.AccessLog.S3BucketName)
}
func resolveElbv1loadBalancerAttributesAccessLogS3BucketPrefix(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.AccessLog == nil {
		return nil
	}
	return resource.Set(c.Name, r.Attributes.AccessLog.S3BucketPrefix)
}
func resolveElbv1loadBalancerAttributesAccessLogEmitInterval(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.AccessLog == nil {
		return nil
	}
	return resource.Set(c.Name, r.Attributes.AccessLog.EmitInterval)
}
func resolveElbv1loadBalancerAttributesConnectionSettingsIdleTimeout(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.ConnectionSettings == nil {
		return nil
	}

	return resource.Set(c.Name, r.Attributes.ConnectionSettings.IdleTimeout)
}
func resolveElbv1loadBalancerAttributesCrossZoneLoadBalancingEnabled(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.CrossZoneLoadBalancing == nil {
		return nil
	}
	return resource.Set(c.Name, r.Attributes.CrossZoneLoadBalancing.Enabled)
}
func resolveElbv1loadBalancerAttributesConnectionDrainingEnabled(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.ConnectionDraining == nil {
		return nil
	}
	return resource.Set(c.Name, r.Attributes.ConnectionDraining.Enabled)
}
func resolveElbv1loadBalancerAttributesConnectionDrainingTimeout(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil && r.Attributes.ConnectionDraining == nil {
		return nil
	}
	return resource.Set(c.Name, r.Attributes.ConnectionDraining.Timeout)
}
func resolveElbv1loadBalancerAttributesAdditionalAttributes(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	if r.Attributes == nil {
		return nil
	}

	response := make(map[string]interface{}, len(r.Attributes.AdditionalAttributes))
	for _, a := range r.Attributes.AdditionalAttributes {
		response[*a.Key] = a.Value
	}
	return resource.Set(c.Name, response)
}
func resolveElbv1loadBalancerInstances(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	response := make([]string, 0, len(r.Instances))
	for _, i := range r.Instances {
		response = append(response, *i.InstanceId)
	}
	return resource.Set(c.Name, response)
}
func fetchElbv1LoadBalancerBackendServerDescriptions(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	r, ok := parent.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	res <- r.BackendServerDescriptions
	return nil
}
func fetchElbv1LoadBalancerListeners(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	r, ok := parent.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	res <- r.ListenerDescriptions
	return nil
}
func fetchElbv1LoadBalancerPoliciesAppCookieStickinessPolicies(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	r, ok := parent.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}

	if r.Policies == nil {
		return nil
	}
	res <- r.Policies.AppCookieStickinessPolicies
	return nil
}
func fetchElbv1LoadBalancerPoliciesLbCookieStickinessPolicies(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	r, ok := parent.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}

	if r.Policies == nil {
		return nil
	}
	res <- r.Policies.LBCookieStickinessPolicies
	return nil
}
func fetchElbv1LoadBalancerPolicies(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	r, ok := parent.Item.(ELBv1LoadBalancerWrapper)
	if !ok {
		return errors.New("not load balancer")
	}
	c := meta.(*client.Client)
	svc := c.Services().ELBv1
	response, err := svc.DescribeLoadBalancerPolicies(ctx, &elbv1.DescribeLoadBalancerPoliciesInput{LoadBalancerName: r.LoadBalancerName}, func(options *elbv1.Options) {
		options.Region = c.Region
	})
	if err != nil {
		return err
	}
	res <- response.PolicyDescriptions
	return nil
}
func resolveElbv1loadBalancerPolicyPolicyAttributeDescriptions(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r, ok := resource.Item.(types.PolicyDescription)
	if !ok {
		return errors.New("not policy description")
	}

	response := make(map[string]interface{}, len(r.PolicyAttributeDescriptions))
	for _, a := range r.PolicyAttributeDescriptions {
		response[*a.AttributeName] = a.AttributeValue
	}
	return resource.Set(c.Name, response)
}
