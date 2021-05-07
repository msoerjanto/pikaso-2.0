package piece

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
	"log"
)

type Artist struct {
	ArtistId  string
	FirstName string
	LastName  string
}

type ArtistRepository interface {
	Store(*Artist) error
	FindAll() ([]*Artist, error)
}

type artistRepository struct {
}

func NewArtistRepository() ArtistRepository {
	return &artistRepository{}
}

func (r *artistRepository) Store(artist *Artist) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000")}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(artist)
	if err != nil {
		log.Fatalf("Got error marshalling new piece item: %s", err)
		return err
	}
	// snippet-end:[dynamodb.go.create_item.assign_struct]

	// snippet-start:[dynamodb.go.create_item.call]
	tableName := "Artists"

	// first check if the item exists
	exists, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ArtistId": {
				S: aws.String(artist.ArtistId),
			},
		},
	})

	if err != nil {
		return errors.New("Got error calling GetItem")
	}
	if exists != nil && exists.Item != nil {
		return errors.New("Artst exists with ArtistId " + artist.ArtistId)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return err
	}

	fmt.Println("Successfully added '" + artist.ArtistId + " to table " + tableName)
	// snippet-end:[dynamodb.go.create_item.call]
	return nil
}

func (r *artistRepository) FindAll() ([]*Artist, error) {
	// snippet-start:[dynamodb.go.load_items.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000")}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.load_items.session]

	tableName := "Artists"
	// snippet-end:[dynamodb.go.scan_items.vars]

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return make([]*Artist, 0), err
	}
	// snippet-end:[dynamodb.go.scan_items.call]

	// snippet-start:[dynamodb.go.scan_items.process]
	var artists []*Artist
	for _, i := range result.Items {
		artist := Artist{}

		err = dynamodbattribute.UnmarshalMap(i, &artist)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
			return artists, err
		}
		artists = append(artists, &artist)
	}
	return artists, nil
}
