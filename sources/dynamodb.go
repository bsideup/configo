package sources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBSource struct {
	Endpoint  string `json:"endpoint"`
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Table     string `json:"table"`
	Key       string `json:"key"`
}

func (dynamoDBSource *DynamoDBSource) Get() (map[string]interface{}, error) {

	config := defaults.Config()

	if dynamoDBSource.AccessKey != "" {
		config = config.WithCredentials(credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     dynamoDBSource.AccessKey,
				SecretAccessKey: dynamoDBSource.SecretKey,
			},
		}))
	}

	if dynamoDBSource.Endpoint != "" {
		config = config.WithEndpoint(dynamoDBSource.Endpoint)
	}

	if dynamoDBSource.Region != "" {
		config = config.WithRegion(dynamoDBSource.Region)
	} else {
		config = config.WithRegion("us-west-1")
	}

	client := dynamodb.New(session.New(config))

	tableName := aws.String(dynamoDBSource.Table)

	err := client.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: tableName})

	if err != nil {
		return nil, err
	}

	key := dynamoDBSource.Key

	response, err := client.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key": {S: aws.String(key)},
		},
		TableName:      tableName,
		ConsistentRead: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = dynamodbattribute.ConvertFromMap(response.Item, &result)

	delete(result, key)

	return result, err
}
