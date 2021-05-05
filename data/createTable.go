package data

// snippet-start:[dynamodb.go.create_table.imports]
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"fmt"
	"log"
)

// snippet-end:[dynamodb.go.create_table.imports]
func CreateTables() {

	//TODO refactor so we only call DoesTableExist once
	// change signature to accept a list of table names to check for
	// returns a map with key table name and value bool representing whether or not it exists
	if tableExists, err := DoesTableExist("Pieces"); err != nil {
		fmt.Println("Problem determining the existence of table")
	} else {
		if !tableExists {
			createPiecesTable()
		}
	}

	if tableExists, err := DoesTableExist("Artists"); err != nil {
		fmt.Println("Problem determining the existence of table")
	} else {
		if !tableExists {
			createArtistsTable()
		}
	}

	if tableExists, err := DoesTableExist("Spaces"); err != nil {
		fmt.Println("Problem determining the existence of table")
	} else {
		if !tableExists {
			createSpacesTable()
		}
	}
}

func createArtistsTable() {
	fmt.Println("Creating Artists table...")
	// snippet-start:[dynamodb.go.create_table.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSession(&aws.Config{

		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000")}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.create_table.session]

	// snippet-start:[dynamodb.go.create_table.call]
	// Create table Movies
	tableName := "Artists"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ArtistId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ArtistId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)
	// snippet-end:[dynamodb.go.create_table.call]
}

func createPiecesTable() {
	fmt.Println("Creating Pieces table...")
	// snippet-start:[dynamodb.go.create_table.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSession(&aws.Config{

		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000")}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.create_table.session]

	// snippet-start:[dynamodb.go.create_table.call]
	// Create table Movies
	tableName := "Pieces"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				// ArtistId+PictureNumber
				AttributeName: aws.String("PieceId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Year"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("PieceId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Year"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)
	// snippet-end:[dynamodb.go.create_table.call]
}

func createSpacesTable() {
	fmt.Println("Creating Spaces table...")
	// snippet-start:[dynamodb.go.create_table.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSession(&aws.Config{

		Region:   aws.String("ap-southeast-1"),
		Endpoint: aws.String("http://localhost:8000")}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	// snippet-end:[dynamodb.go.create_table.session]

	// snippet-start:[dynamodb.go.create_table.call]
	// Create table Movies
	tableName := "Spaces"

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Location"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SpaceNumber"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Location"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("SpaceNumber"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	fmt.Println("Created the table", tableName)
	// snippet-end:[dynamodb.go.create_table.call]
}
