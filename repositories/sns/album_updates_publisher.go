package repositories

import (
	"example/golang-api/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/zap"
)

type AlbumUpdatesPublisher struct {
	snsClient           *sns.SNS
	albumUpdateTopicARN *string
	slog                *zap.SugaredLogger
}

func NewAlbumUpdatesPublisher(config *config.Config, slogger *zap.SugaredLogger) *AlbumUpdatesPublisher {

	albumUpdateTopicARN := config.Sns.AlbumUpdateTopicArn

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      &config.Sns.Region,
		Credentials: credentials.NewStaticCredentials(
      config.Sns.Credentials.Id,
      config.Sns.Credentials.Secret,
      config.Sns.Credentials.Token,
    ),
		Endpoint:    &config.Sqs.ServerEndpoint,
	}))

	endpoint := config.Sns.ServerEndpoint

	snsClient := sns.New(sess, &aws.Config{
		Endpoint: &endpoint,
	})

	aup := AlbumUpdatesPublisher{
		snsClient:           snsClient,
		albumUpdateTopicARN: &albumUpdateTopicARN,
		slog:                slogger,
	}

	return &aup
}

func (aup AlbumUpdatesPublisher) PublishUpdatedAlbum(message string) (string, error) {

	result, err := aup.snsClient.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: aup.albumUpdateTopicARN,
	})
	if err != nil {
		aup.slog.Warnf("sns pub failed %v", err)
		return "", err
	}

	return *result.MessageId, nil
}
