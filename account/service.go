package account

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	ddb *dynamodb.Client
}

type Account struct {
	Email          string  `json:"email"`
	HashedPassword string  `json:"-"`
	Balance        float64 `json:"balance"`
}

type ActionType string

const (
	Pay     ActionType = "pay"
	Receive ActionType = "receive"
)

type Transaction struct {
	Id          string     `json:"id"`
	Amount      float64    `json:"amount"`
	Timestamp   time.Time  `json:"timestamp"`
	Email       string     `json:"email"`
	Participant string     `json:"participant"`
	Action      ActionType `json:"action"`
}

func NewFromClient(dynamodbClient *dynamodb.Client) (*AccountService, error) {
	return &AccountService{
		ddb: dynamodbClient,
	}, nil
}

func (as AccountService) Register(ctx context.Context, email, password string) (Account, error) {
	account := Account{}
	account.Email = email

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}
	account.HashedPassword = string(hashedBytes)
	account.Balance = 0.0

	putItemInput := dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("ACCOUNTSTABLE")),
		Item: map[string]types.AttributeValue{
			"PK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", account.Email)},
			"SK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", account.Email)},
			"email":    &types.AttributeValueMemberS{Value: account.Email},
			"balance":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", account.Balance)},
			"password": &types.AttributeValueMemberS{Value: string(hashedBytes)},
		},
		ConditionExpression: aws.String("attribute_not_exists(SK)"),
	}

	_, err = as.ddb.PutItem(ctx, &putItemInput)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func (as AccountService) Login(ctx context.Context, email, password string) (string, error) {

	getItemInput := dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("ACCOUNTSTABLE")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", email)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", email)},
		},
	}

	getItemResponse, err := as.ddb.GetItem(ctx, &getItemInput)
	if err != nil {
		return "", err
	}

	if getItemResponse.Item == nil {
		err := fmt.Errorf("account not found: %s", email)
		return "", err
	}
	hashedPassword := getItemResponse.Item["password"].(*types.AttributeValueMemberS).Value
	savedEmail := getItemResponse.Item["email"].(*types.AttributeValueMemberS).Value
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", err
	}

	claims := &jwt.StandardClaims{
		Subject:   savedEmail,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("SECRET")))

	return ss, nil
}

func (as AccountService) GetByEmail(ctx context.Context, email string) (Account, error) {
	getItemInput := dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("ACCOUNTSTABLE")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", email)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", email)},
		},
	}

	getItemResponse, err := as.ddb.GetItem(ctx, &getItemInput)
	if err != nil {
		return Account{}, err
	}

	if getItemResponse.Item == nil {
		err := fmt.Errorf("account not found: %s", email)
		return Account{}, err
	}

	balance, err := strconv.ParseFloat(getItemResponse.Item["balance"].(*types.AttributeValueMemberN).Value, 64)
	if err != nil {
		return Account{}, err
	}

	account := Account{
		Email:   getItemResponse.Item["email"].(*types.AttributeValueMemberS).Value,
		Balance: balance,
	}

	return account, nil
}

func (as AccountService) GetTransactions(ctx context.Context, email string) ([]Transaction, error) {
	transactions := make([]Transaction, 0)
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("ACCOUNTSTABLE")),
		KeyConditionExpression: aws.String(":p = #p AND begins_with(:s, #s)"),
		ExpressionAttributeNames: map[string]string{
			":p": "PK",
			":s": "SK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			"#p": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", email)},
			"#s": &types.AttributeValueMemberS{Value: "T#"},
		},
	}

	queryResponse, err := as.ddb.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}

	// TODO setup pagination or loop for large transaction list

	for _, t := range queryResponse.Items {
		amount, err := strconv.ParseFloat(t["amount"].(*types.AttributeValueMemberN).Value, 64)
		if err != nil {
			return nil, err
		}
		timestamp, err := time.Parse(time.RFC3339, t["timestamp"].(*types.AttributeValueMemberS).Value)
		if err != nil {
			return nil, err
		}
		transaction := Transaction{
			Id:          t["id"].(*types.AttributeValueMemberS).Value,
			Email:       t["email"].(*types.AttributeValueMemberS).Value,
			Amount:      amount,
			Participant: t["participant"].(*types.AttributeValueMemberS).Value,
			Action:      ActionType(t["action"].(*types.AttributeValueMemberS).Value),
			Timestamp:   timestamp,
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
