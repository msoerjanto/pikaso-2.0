package piece

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
	"log"
	"strconv"
)

type Piece struct {
	PieceId  string // ArtistId+PictureNumber
	Year     int
	Title    string
	Media    string
	Length   int
	Height   int
	Page     int
	ImageUrl string
}

type Repository interface {
	Store(piece *Piece) error
	FindAll() ([]*Piece, error)
}

type pieceRepository struct {
}

func NewPieceRepository() Repository {
	return &pieceRepository{}
}

func (r *pieceRepository) Store(piece *Piece) error {
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

	av, err := dynamodbattribute.MarshalMap(piece)
	if err != nil {
		log.Fatalf("Got error marshalling new piece item: %s", err)
		return err
	}
	// snippet-end:[dynamodb.go.create_item.assign_struct]

	// snippet-start:[dynamodb.go.create_item.call]
	// Create item in table Movies
	tableName := "Pieces"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return err
	}

	year := strconv.Itoa(piece.Year)

	fmt.Println("Successfully added '" + piece.Title + "' (" + year + ") to table " + tableName)
	// snippet-end:[dynamodb.go.create_item.call]
	return nil
}

func (r *pieceRepository) FindAll() ([]*Piece, error) {
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

	tableName := "Pieces"
	// snippet-end:[dynamodb.go.scan_items.vars]

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
		return make([]*Piece, 0), err
	}
	// snippet-end:[dynamodb.go.scan_items.call]

	// snippet-start:[dynamodb.go.scan_items.process]
	var pieces []*Piece
	for _, i := range result.Items {
		piece := Piece{}

		err = dynamodbattribute.UnmarshalMap(i, &piece)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
			return pieces, err
		}
		pieces = append(pieces, &piece)
	}
	return pieces, nil
}
