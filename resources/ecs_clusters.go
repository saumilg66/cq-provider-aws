package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/cloudquery/cq-provider-aws/client"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func EcsClusters() *schema.Table {
	return &schema.Table{
		Name:         "aws_ecs_clusters",
		Resolver:     fetchEcsClusters,
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
				Name: "active_services_count",
				Type: schema.TypeInt,
			},
			{
				Name: "attachments_status",
				Type: schema.TypeString,
			},
			{
				Name: "capacity_providers",
				Type: schema.TypeStringArray,
			},
			{
				Name: "cluster_arn",
				Type: schema.TypeString,
			},
			{
				Name: "cluster_name",
				Type: schema.TypeString,
			},
			{
				Name:     "execute_config_kms_key_id",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Configuration.ExecuteCommandConfiguration.KmsKeyId"),
			},
			{
				Name:     "execute_config_logs_cloud_watch_encryption_enabled",
				Type:     schema.TypeBool,
				Resolver: schema.PathResolver("Configuration.ExecuteCommandConfiguration.LogConfiguration.CloudWatchEncryptionEnabled"),
			},
			{
				Name:     "execute_config_log_cloud_watch_log_group_name",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Configuration.ExecuteCommandConfiguration.LogConfiguration.CloudWatchLogGroupName"),
			},
			{
				Name:     "execute_config_log_s3_bucket_name",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Configuration.ExecuteCommandConfiguration.LogConfiguration.S3BucketName"),
			},
			{
				Name:     "execute_config_log_s3_encryption_enabled",
				Type:     schema.TypeBool,
				Resolver: schema.PathResolver("Configuration.ExecuteCommandConfiguration.LogConfiguration.S3EncryptionEnabled"),
			},
			{
				Name:     "execute_config_log_s3_key_prefix",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Configuration.ExecuteCommandConfiguration.LogConfiguration.S3KeyPrefix"),
			},
			{
				Name:     "execute_config_logging",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Configuration.ExecuteCommandConfiguration.Logging"),
			},
			{
				Name: "pending_tasks_count",
				Type: schema.TypeInt,
			},
			{
				Name: "registered_container_instances_count",
				Type: schema.TypeInt,
			},
			{
				Name: "running_tasks_count",
				Type: schema.TypeInt,
			},
			{
				Name:     "settings",
				Type:     schema.TypeJSON,
				Resolver: resolveEcsClusterSettings,
			},
			{
				Name:     "statistics",
				Type:     schema.TypeJSON,
				Resolver: resolveEcsClusterStatistics,
			},
			{
				Name: "status",
				Type: schema.TypeString,
			},
			{
				Name:     "tags",
				Type:     schema.TypeJSON,
				Resolver: resolveEcsClusterTags,
			},
		},
		Relations: []*schema.Table{
			{
				Name:     "aws_ecs_cluster_attachments",
				Resolver: fetchEcsClusterAttachments,
				Columns: []schema.Column{
					{
						Name:     "cluster_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name:     "attachment_id",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("Id"),
					},
					{
						Name: "status",
						Type: schema.TypeString,
					},
					{
						Name: "type",
						Type: schema.TypeString,
					},
					{
						Name:     "details",
						Type:     schema.TypeJSON,
						Resolver: resolveEcsClusterAttachmentDetails,
					},
				},
			},
			{
				Name:     "aws_ecs_cluster_default_capacity_provider_strategies",
				Resolver: fetchEcsClusterDefaultCapacityProviderStrategies,
				Columns: []schema.Column{
					{
						Name:     "cluster_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name: "capacity_provider",
						Type: schema.TypeString,
					},
					{
						Name: "base",
						Type: schema.TypeInt,
					},
					{
						Name: "weight",
						Type: schema.TypeInt,
					},
				},
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================
func fetchEcsClusters(ctx context.Context, meta schema.ClientMeta, _ *schema.Resource, res chan interface{}) error {
	var config ecs.ListClustersInput
	region := meta.(*client.Client).Region
	svc := meta.(*client.Client).Services().ECS
	for {
		listClustersOutput, err := svc.ListClusters(ctx, &config, func(o *ecs.Options) {
			o.Region = region
		})
		if err != nil {
			return err
		}
		describeClusterOutput, err := svc.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: listClustersOutput.ClusterArns}, func(o *ecs.Options) {
			o.Region = region
		})
		if err != nil {
			return err
		}
		res <- describeClusterOutput.Clusters

		if listClustersOutput.NextToken == nil {
			break
		}
		config.NextToken = listClustersOutput.NextToken
	}
	return nil
}
func resolveEcsClusterSettings(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, _ schema.Column) error {
	cluster := resource.Item.(types.Cluster)
	settings := make(map[string]*string)
	for _, s := range cluster.Settings {
		settings[string(s.Name)] = s.Value
	}
	return resource.Set("settings", settings)
}
func resolveEcsClusterStatistics(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, _ schema.Column) error {
	cluster := resource.Item.(types.Cluster)
	stats := make(map[string]*string)
	for _, s := range cluster.Statistics {
		stats[*s.Name] = s.Value
	}
	return resource.Set("statistics", stats)
}
func resolveEcsClusterTags(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, _ schema.Column) error {
	cluster := resource.Item.(types.Cluster)
	stats := make(map[string]*string)
	for _, s := range cluster.Tags {
		stats[*s.Key] = s.Value
	}
	return resource.Set("tags", stats)
}
func fetchEcsClusterAttachments(_ context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	cluster := parent.Item.(types.Cluster)
	res <- cluster.Attachments
	return nil
}
func resolveEcsClusterAttachmentDetails(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, _ schema.Column) error {
	attachment := resource.Item.(types.Attachment)
	details := make(map[string]*string)
	for _, s := range attachment.Details {
		details[*s.Name] = s.Value
	}
	return resource.Set("details", details)
}
func fetchEcsClusterDefaultCapacityProviderStrategies(_ context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	cluster := parent.Item.(types.Cluster)
	res <- cluster.DefaultCapacityProviderStrategy
	return nil
}
