package space

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"log"
)

type Space struct {
	Location    string
	SpaceNumber int
	PieceId     string // ArtistId + PictureNumber
}

type SpaceRepository interface {
	Store(space *Space) error
	FindAll() ([]*Space, error)
}

type spaceRepository struct {
}

func NewSpaceRepository() SpaceRepository {
	return &spaceRepository{}
}

func (s *spaceRepository) Store(space *Space) error {
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
		return err
	}
	// snippet-end:[dynamodb.go.create_item.assign_struct]

	// snippet-start:[dynamodb.go.create_item.call]
	// Create item in table Movies
	tableName := "Spaces"

	// first check if the item exists
	exists, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Location": {
				S: aws.String(space.Location),
			},
			"SpaceNumber": {
				N: aws.String(strconv.Itoa(space.SpaceNumber)),
			},
		},
	})

	if err != nil {
		return errors.New("Got error calling GetItem")
	}
	if exists != nil && exists.Item != nil {
		return errors.New("Space already exists")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully added '%s %d' to table %s \n", space.Location, space.SpaceNumber, tableName)
	return nil
	// snippet-end:[dynamodb.go.create_item.call]
}

func (s *spaceRepository) FindAll() ([]*Space, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000")}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.load_items.session]

	tableName := "Spaces"
	// snippet-end:[dynamodb.go.scan_items.vars]

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return make([]*Space, 0), err
	}
	// snippet-end:[dynamodb.go.scan_items.call]

	// snippet-start:[dynamodb.go.scan_items.process]
	var spaces []*Space
	for _, i := range result.Items {
		space := Space{}

		err = dynamodbattribute.UnmarshalMap(i, &space)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
			return spaces, err
		}
		spaces = append(spaces, &space)
	}
	return spaces, nil
}
