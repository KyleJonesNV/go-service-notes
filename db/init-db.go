package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/KyleJonesNV/go-service-notes/pkg/notes"
)

type userTopic struct {

}

func main() {
	ctx := context.Background()

	insertUser, err := readUser("user.json")
	if err != nil {
		log.Fatal(err)
	}
	if insertUser == nil {
		log.Fatalf("insert user is nil")
	}

	topicsToInsert, err := readTopics("topics.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserting: ", insertUser.Email)
	user, err := notes.InsertUser(ctx, *insertUser)
	if err != nil {
		log.Fatal(err)
	}
	if user == nil {
		log.Fatalf("user is nil")
	}
	if user.ID == "" {
		log.Fatalf("user id nil")
	}

	for _, topic := range topicsToInsert {
		fmt.Println("inserting: ", topic.Title)
		err = notes.InsertTopic(ctx, user.ID, topic.Title)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readTopics(fileName string) ([]notes.Topic, error) {
	input, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read file, %w", err)
	}

	var topics = []notes.Topic{}
	err = json.Unmarshal(input, &topics)
	if err != nil {
		return nil, fmt.Errorf("unmarshal, %w", err)
	}

	return topics, nil
}

func readUser(fileName string) (*notes.UserInsert, error) {
	input, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read file, %w", err)
	}

	var user = notes.UserInsert{}
	err = json.Unmarshal(input, &user)
	if err != nil {
		return nil, fmt.Errorf("unmarshal, %w", err)
	}

	return &user, nil
}

// func insertMovie(ctx context.Context, cfg aws.Config, movie movies.Movie) error {
// 	item, err := attributevalue.MarshalMap(movie)
// 	if err != nil {
// 		return fmt.Errorf("dynamo marshal map, %w", err)
// 	}

// 	svc := dynamodb.NewFromConfig(cfg)

// 	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
// 		TableName: aws.String(movies.TableName),
// 		Item:      item,
// 	})

// 	if err != nil {
// 		return fmt.Errorf("dynamo put item, %w", err)
// 	}

// 	return nil
// }
