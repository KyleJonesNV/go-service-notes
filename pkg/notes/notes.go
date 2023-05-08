package notes

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofrs/uuid"
)

const TableName = "go-service-notes"

type User struct {
	ID string
	Email string
	Name string
	Surname string
}

type UserInsert struct {
	Email string
	Name string
	Surname string
}

type Topic struct {
	Title string
	Notes []Note
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Note struct {
	Title string
	Content string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	userPrefix = "user"
	topicPrefix = "topic"
)

const (
	pk = "PK"
	sk = "SK"
)

type KeyValue struct {
	Key string
	Value string
}

type DBKey struct {
	Hash KeyValue
	Sort KeyValue
}

func userKey(userEmail string) DBKey {
	return DBKey{
		Hash: KeyValue{
			Key: pk,
			Value: userPrefix,
		},
		Sort: KeyValue{
			Key: sk,
			Value: userEmail,
		},
	}
}

func topicKey(userID string, title string) DBKey {
	return DBKey{
		Hash: KeyValue{
			Key: pk,
			Value: fmt.Sprintf("%s#%s", topicPrefix, userID),
		},
		Sort: KeyValue{
			Key: sk,
			Value: title,
		},
	}
}


func GetAllForUser(ctx context.Context, UserID string) ([]Topic, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load default config, %w", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	keyCond := expression.Key(pk).Equal(expression.Value(topicKey(UserID, "").Hash.Value))
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("expression builder: %w", err)
	}

	queryInput := dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName: aws.String(TableName),
	}

	resp, err := svc.Query(ctx, &queryInput)

	if err != nil {
		return nil, fmt.Errorf("query, %w", err)
	}

	var topics = []Topic{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &topics)
	if err != nil {
		return nil, fmt.Errorf("unmarshal list of maps, %w", err)
	}

	return topics, nil
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load default config, %w", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	keyCond := expression.Key(pk).Equal(expression.Value(userKey("").Hash.Value))
	expression.Key(sk).Equal(expression.Value(email))
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("expression builder: %w", err)
	}

	queryInput := dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName: aws.String(TableName),
	}

	resp, err := svc.Query(ctx, &queryInput)

	if err != nil {
		return nil, fmt.Errorf("query, %w", err)
	}

	var users = []User{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &users)
	if err != nil {
		return nil, fmt.Errorf("unmarshal list of maps, %w", err)
	}
	if len(users) == 0 {
		return nil, nil
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("more than 1 user with email %q, an error has occured in dynamo setup", email)
	}

	return &users[0], nil
}

func GetUserTopicByTitle(ctx context.Context, userID, title string) (*Topic, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load default config, %w", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	keyCond := expression.Key(pk).Equal(expression.Value(topicKey(userID, userID).Hash.Value))
	keyCond = keyCond.And(expression.Key(sk).Equal(expression.Value(topicKey(userID, title).Sort.Value)))
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	expr, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("expression builder: %w", err)
	}

	queryInput := dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		TableName: aws.String(TableName),
	}

	resp, err := svc.Query(ctx, &queryInput)

	if err != nil {
		return nil, fmt.Errorf("query, %w", err)
	}

	var topic = []Topic{}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, &topic)
	if err != nil {
		return nil, fmt.Errorf("unmarshal list of maps, %w", err)
	}
	if len(topic) == 0 {
		return nil, nil
	}

	if len(topic) > 1 {
		return nil, fmt.Errorf("more than 1 topic found with the same title: %q userID: %q, an error has occured in dynamo setup, topics: %q, %q", title, userID, topic[0].Title, topic[1].Title)
	}

	return &topic[0], nil
}

func InsertUser(ctx context.Context, userInsert UserInsert) (*User, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load default config, %w", err)
	}

	foundUser, err := GetUserByEmail(ctx, userInsert.Email)
	if err != nil {
		return nil, fmt.Errorf("get user by email, %w", err)
	}

	if foundUser != nil {
		return foundUser, nil
	}

	userID := uuid.Must(uuid.NewV4()).String()

	user := User{
		ID: userID,
		Name: userInsert.Name,
		Surname: userInsert.Surname,
		Email: userInsert.Email,
	}

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, fmt.Errorf("dynamo marshal map, %w", err)
	}

	key := userKey(userInsert.Email)

	hashValue, err := attributevalue.Marshal(aws.String(key.Hash.Value))
	if err != nil {
		return nil,  fmt.Errorf("marshal hash %q: %w", key.Hash.Value, err)
	}
	sortValue, err := attributevalue.Marshal(aws.String(key.Sort.Value))
	if err != nil {
		return nil,  fmt.Errorf("marshal sort %q: %w", key.Sort.Value, err)
	}

	item[key.Hash.Key] = hashValue
	item[key.Sort.Key] = sortValue

	item[key.Hash.Key] = hashValue

	svc := dynamodb.NewFromConfig(cfg)

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})

	if err != nil {
		return nil, fmt.Errorf("dynamo put item, %w", err)
	}

	return &user, nil
}

func InsertTopic(ctx context.Context, userID string, title string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("load default config, %w", err)
	}

	topic := Topic{
		Title:	title,
		CreatedAt: time.Now().UTC(),
	}

	item, err := attributevalue.MarshalMap(topic)
	if err != nil {
		return fmt.Errorf("dynamo marshal map, %w", err)
	}

	key := topicKey(userID, title)

	hashValue, err := attributevalue.Marshal(aws.String(key.Hash.Value))
	if err != nil {
		return fmt.Errorf("marshal hash %q: %w", key.Hash.Value, err)
	}
	sortValue, err := attributevalue.Marshal(aws.String(key.Sort.Value))
	if err != nil {
		return fmt.Errorf("marshal sort %q: %w", key.Sort.Value, err)
	}

	item[key.Hash.Key] = hashValue
	item[key.Sort.Key] = sortValue

	svc := dynamodb.NewFromConfig(cfg)

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("dynamo put item, %w", err)
	}

	return nil
}

func DeleteTopic(ctx context.Context, userID string, title string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("load default config, %w", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	_, err = svc.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key:       getTopicKey(userID, title),
	})

	if err != nil {
		return fmt.Errorf("dynamo put item, %w", err)
	}

	return nil
}

func InsertNote(ctx context.Context, userID string, title string, note Note) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("load default config, %w", err)
	}

	topic, err := GetUserTopicByTitle(ctx, userID, title)
	if err != nil {
		return fmt.Errorf("get user topic by title, %w", err)
	}
	if topic == nil {
		return fmt.Errorf("unknown topic %q, for userID %q", title, userID)
	}

	topic.Notes = append(topic.Notes, note)

	item, err := attributevalue.MarshalMap(*topic)
	if err != nil {
		return fmt.Errorf("dynamo marshal map, %w", err)
	}

	key := topicKey(userID, title)

	hashValue, err := attributevalue.Marshal(aws.String(key.Hash.Value))
	if err != nil {
		return fmt.Errorf("marshal hash %q: %w", key.Hash.Value, err)
	}
	sortValue, err := attributevalue.Marshal(aws.String(key.Sort.Value))
	if err != nil {
		return fmt.Errorf("marshal sort %q: %w", key.Sort.Value, err)
	}

	item[key.Hash.Key] = hashValue
	item[key.Sort.Key] = sortValue

	svc := dynamodb.NewFromConfig(cfg)

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("dynamo put item, %w", err)
	}

	return nil
}

func DeleteNote(ctx context.Context, userID string, title string, NoteTitle string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("load default config, %w", err)
	}

	topic, err := GetUserTopicByTitle(ctx, userID, title)
	if err != nil {
		return fmt.Errorf("get user topic by title, %w", err)
	}
	if topic == nil {
		return fmt.Errorf("unknown topic %q, for userID %q", title, userID)
	}

	for s, note := range topic.Notes{
		if note.Title == NoteTitle {
			topic.Notes = append(topic.Notes[:s], topic.Notes[s+1:]...)	
		}		
	}

	item, err := attributevalue.MarshalMap(*topic)
	if err != nil {
		return fmt.Errorf("dynamo marshal map, %w", err)
	}

	key := topicKey(userID, title)

	hashValue, err := attributevalue.Marshal(aws.String(key.Hash.Value))
	if err != nil {
		return fmt.Errorf("marshal hash %q: %w", key.Hash.Value, err)
	}
	sortValue, err := attributevalue.Marshal(aws.String(key.Sort.Value))
	if err != nil {
		return fmt.Errorf("marshal sort %q: %w", key.Sort.Value, err)
	}

	item[key.Hash.Key] = hashValue
	item[key.Sort.Key] = sortValue

	svc := dynamodb.NewFromConfig(cfg)

	_, err = svc.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("dynamo put item, %w", err)
	}

	return nil
}

func getTopicKey(userID string, title string) map[string]types.AttributeValue {
	key := topicKey(userID, title)
	hash, err := attributevalue.Marshal(key.Hash.Value)
	if err != nil {
		panic(err)
	}
	sort, err := attributevalue.Marshal(key.Sort.Value)
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{pk: hash, sk: sort}
}

func getTopicNoteKey(userID string, title string) map[string]types.AttributeValue {
	key := topicKey(userID, title)
	hash, err := attributevalue.Marshal(key.Hash.Value)
	if err != nil {
		panic(err)
	}
	sort, err := attributevalue.Marshal(key.Sort.Value)
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{pk: hash, sk: sort}
}
