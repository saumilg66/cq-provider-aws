package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/cloudquery/cq-provider-aws/client"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func CloudtrailTrails() *schema.Table {
	return &schema.Table{
		Name:                 "aws_cloudtrail_trails",
		Resolver:             fetchCloudtrailTrails,
		Multiplex:            client.AccountRegionMultiplex,
		IgnoreError:          client.IgnoreAccessDeniedServiceDisabled,
		DeleteFilter:         client.DeleteAccountRegionFilter,
		PostResourceResolver: postCloudtrailTrailResolver,
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
				Name:     "cloudwatch_logs_log_group_name",
				Type:     schema.TypeString,
				Resolver: resolveCloudtrailTrailCloudwatchLogsLogGroupName,
			},
			{
				Name: "is_logging",
				Type: schema.TypeBool,
			},
			{
				Name: "latest_cloud_watch_logs_delivery_error",
				Type: schema.TypeString,
			},
			{
				Name: "latest_cloud_watch_logs_delivery_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "latest_delivery_attempt_succeeded",
				Type: schema.TypeString,
			},
			{
				Name: "latest_delivery_attempt_time",
				Type: schema.TypeString,
			},
			{
				Name: "latest_delivery_error",
				Type: schema.TypeString,
			},
			{
				Name: "latest_delivery_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "latest_digest_delivery_error",
				Type: schema.TypeString,
			},
			{
				Name: "latest_digest_delivery_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "latest_notification_attempt_succeeded",
				Type: schema.TypeString,
			},
			{
				Name: "latest_notification_attempt_time",
				Type: schema.TypeString,
			},
			{
				Name: "latest_notification_error",
				Type: schema.TypeString,
			},
			{
				Name: "latest_notification_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "start_logging_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "stop_logging_time",
				Type: schema.TypeTimestamp,
			},
			{
				Name: "time_logging_started",
				Type: schema.TypeString,
			},
			{
				Name: "time_logging_stopped",
				Type: schema.TypeString,
			},
			{
				Name: "cloud_watch_logs_log_group_arn",
				Type: schema.TypeString,
			},
			{
				Name: "cloud_watch_logs_role_arn",
				Type: schema.TypeString,
			},
			{
				Name: "has_custom_event_selectors",
				Type: schema.TypeBool,
			},
			{
				Name: "has_insight_selectors",
				Type: schema.TypeBool,
			},
			{
				Name: "home_region",
				Type: schema.TypeString,
			},
			{
				Name: "include_global_service_events",
				Type: schema.TypeBool,
			},
			{
				Name: "is_multi_region_trail",
				Type: schema.TypeBool,
			},
			{
				Name: "is_organization_trail",
				Type: schema.TypeBool,
			},
			{
				Name: "kms_key_id",
				Type: schema.TypeString,
			},
			{
				Name: "log_file_validation_enabled",
				Type: schema.TypeBool,
			},
			{
				Name: "name",
				Type: schema.TypeString,
			},
			{
				Name: "s3_bucket_name",
				Type: schema.TypeString,
			},
			{
				Name: "s3_key_prefix",
				Type: schema.TypeString,
			},
			{
				Name:     "sns_topic_arn",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("SnsTopicARN"),
			},
			{
				Name: "sns_topic_name",
				Type: schema.TypeString,
			},
			{
				Name:     "trail_arn",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("TrailARN"),
			},
		},
		Relations: []*schema.Table{
			{
				Name:     "aws_cloudtrail_trail_event_selectors",
				Resolver: fetchCloudtrailTrailEventSelectors,
				Columns: []schema.Column{
					{
						Name:     "trail_id",
						Type:     schema.TypeUUID,
						Resolver: schema.ParentIdResolver,
					},
					{
						Name: "exclude_management_event_sources",
						Type: schema.TypeStringArray,
					},
					{
						Name: "include_management_events",
						Type: schema.TypeBool,
					},
					{
						Name: "read_write_type",
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
func fetchCloudtrailTrails(ctx context.Context, meta schema.ClientMeta, _ *schema.Resource, res chan interface{}) error {
	c := meta.(*client.Client)
	svc := c.Services().Cloudtrail
	response, err := svc.DescribeTrails(ctx, nil, func(options *cloudtrail.Options) {
		options.Region = c.Region
	})
	if err != nil {
		return err
	}
	res <- response.TrailList
	return nil
}
func postCloudtrailTrailResolver(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
	c := meta.(*client.Client)
	svc := c.Services().Cloudtrail
	r := resource.Item.(types.Trail)
	response, err := svc.GetTrailStatus(ctx,
		&cloudtrail.GetTrailStatusInput{Name: r.TrailARN}, func(o *cloudtrail.Options) {
			o.Region = c.Region
		})
	if err != nil {
		return err
	}
	if err := resource.Set("is_logging", response.IsLogging); err != nil {
		return err
	}
	if err := resource.Set("latest_cloud_watch_logs_delivery_error", response.LatestCloudWatchLogsDeliveryError); err != nil {
		return err
	}
	if err := resource.Set("latest_cloud_watch_logs_delivery_time", response.LatestCloudWatchLogsDeliveryTime); err != nil {
		return err
	}
	if err := resource.Set("latest_delivery_attempt_succeeded", response.LatestDeliveryAttemptSucceeded); err != nil {
		return err
	}
	if err := resource.Set("latest_delivery_attempt_time", response.LatestDeliveryAttemptTime); err != nil {
		return err
	}
	if err := resource.Set("latest_delivery_error", response.LatestDeliveryError); err != nil {
		return err
	}
	if err := resource.Set("latest_delivery_time", response.LatestDeliveryTime); err != nil {
		return err
	}
	if err := resource.Set("latest_digest_delivery_error", response.LatestDigestDeliveryError); err != nil {
		return err
	}
	if err := resource.Set("latest_digest_delivery_time", response.LatestDigestDeliveryTime); err != nil {
		return err
	}
	if err := resource.Set("latest_notification_attempt_succeeded", response.LatestNotificationAttemptSucceeded); err != nil {
		return err
	}
	if err := resource.Set("latest_notification_attempt_time", response.LatestNotificationAttemptTime); err != nil {
		return err
	}
	if err := resource.Set("latest_notification_error", response.LatestNotificationError); err != nil {
		return err
	}
	if err := resource.Set("latest_notification_time", response.LatestNotificationTime); err != nil {
		return err
	}
	if err := resource.Set("start_logging_time", response.StartLoggingTime); err != nil {
		return err
	}
	if err := resource.Set("stop_logging_time", response.StopLoggingTime); err != nil {
		return err
	}
	if err := resource.Set("time_logging_started", response.TimeLoggingStarted); err != nil {
		return err
	}
	if err := resource.Set("time_logging_stopped", response.TimeLoggingStopped); err != nil {
		return err
	}
	return nil
}
func resolveCloudtrailTrailCloudwatchLogsLogGroupName(_ context.Context, meta schema.ClientMeta, resource *schema.Resource, _ schema.Column) error {
	groupName := ""
	log := meta.(*client.Client).Logger()
	trail := resource.Item.(types.Trail)
	if trail.CloudWatchLogsLogGroupArn != nil {
		matches := client.GroupNameRegex.FindStringSubmatch(*trail.CloudWatchLogsLogGroupArn)
		if len(matches) < 2 {
			log.Warn("CloudWatchLogsLogGroupARN doesn't fit standard regex", "arn", *trail.CloudWatchLogsLogGroupArn)
		} else {
			groupName = matches[1]
		}
	} else {
		log.Info("CloudWatchLogsLogGroupARN is empty")
	}

	return resource.Set("cloudwatch_logs_log_group_name", groupName)
}

func fetchCloudtrailTrailEventSelectors(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan interface{}) error {
	r := parent.Item.(types.Trail)
	c := meta.(*client.Client)
	svc := c.Services().Cloudtrail
	response, err := svc.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{TrailName: r.TrailARN}, func(options *cloudtrail.Options) {
		options.Region = c.Region
	})
	if err != nil {
		return err
	}
	res <- response.EventSelectors
	return nil
}
