package e2e

import (
	"context"
	"testing"

	"github.com/andyday/depot/dynamo"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/suite"
)

type DynamoSuite struct {
	Suite
}

func TestDynamoSuite(t *testing.T) {
	suite.Run(t, new(DynamoSuite))
}

func (s *DynamoSuite) SetupSuite() {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithClientLogMode(
		aws.LogRequestWithBody|
			aws.LogResponseWithBody|
			aws.LogRequestEventMessage|
			aws.LogResponseEventMessage,
	))
	s.NoError(err)
	s.db, err = dynamo.NewDatabase(cfg)
	s.NoError(err)
}

func (s *DynamoSuite) TestPut() {
	TestPut(&s.Suite)
}

func (s *DynamoSuite) TestCreate() {
	TestCreate(&s.Suite)
}

func (s *DynamoSuite) TestUpdate() {
	TestUpdate(&s.Suite)
}
