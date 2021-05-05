package space

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"log"
)

// snippet-end:[dynamodb.go.create_item.imports]

// snippet-start:[dynamodb.go.create_item.struct]
// Create struct to hold info about new item
type Space struct {
	Location    string
	SpaceNumber int
	PieceId     string // ArtistId + PictureNumber
}

// snippet-end:[dynamodb.go.create_item.struct]

func CreateNewSpace(space Space) {
	// snippet-start:[dynamodb.go.create_item.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000")}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.create_item.session]

	av, err := dynamodbattribute.MarshalMap(space)
	if err != nil {
		log.Fatalf("Got error marshalling new piece item: %s", err)
	}
	// snippet-end:[dynamodb.go.create_item.assign_struct]

	// snippet-start:[dynamodb.go.create_item.call]
	// Create item in table Movies
	tableName := "Spaces"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Printf("Successfully added '%s %d' to table %s \n", space.Location, space.SpaceNumber, tableName)

	// snippet-end:[dynamodb.go.create_item.call]
}
