package resources

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/cloudquery/cq-provider-aws/client"
	"github.com/cloudquery/cq-provider-aws/client/mocks"
	"github.com/golang/mock/gomock"
)

func buildEc2RegionalConfig(t *testing.T, ctrl *gomock.Controller) client.Services {
	m := mocks.NewMockEc2Client(ctrl)
	m.EXPECT().GetEbsDefaultKmsKeyId(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ec2.GetEbsDefaultKmsKeyIdOutput{KmsKeyId: aws.String("some/key/id")}, nil)
	m.EXPECT().GetEbsEncryptionByDefault(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ec2.GetEbsEncryptionByDefaultOutput{EbsEncryptionByDefault: true}, nil)

	return client.Services{
		EC2: m,
	}
}

func TestEc2RegionalConfig(t *testing.T) {
	awsTestHelper(t, Ec2RegionalConfig(), buildEc2RegionalConfig)
}
