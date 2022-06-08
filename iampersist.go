package iampersist

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
)

type CreatePersistenceInput struct {
	AccessKey    string
	SecretKey    string
	SessionToken string
	Username     string
}

// SetAccessKey - set CreatePersistenceInput.AccessKey value
func (c *CreatePersistenceInput) SetAccessKey(v string) *CreatePersistenceInput {
	c.AccessKey = v
	return c
}

// SetSecretKey - set CreatePersistenceInput.SecretKey value
func (c *CreatePersistenceInput) SetSecretKey(v string) *CreatePersistenceInput {
	c.SecretKey = v
	return c
}

// SetSessionToken - set CreatePersistenceInput.SessionToken value
func (c *CreatePersistenceInput) SetSessionToken(v string) *CreatePersistenceInput {
	c.SessionToken = v
	return c
}

// SetUsername - set CreatePersistenceInput.Username value
func (c *CreatePersistenceInput) SetUsername(v string) *CreatePersistenceInput {
	c.Username = v
	return c
}

// CreatePersistence - create the aws iam persistence and return the credentials
// CreatePersistenceInput
// 	* AccessKey - required
// 	* SecretKey - required
// 	* SessionToken 
//  * Username - when username unset this function search for aws user and create new access key  randomly,
//   but if user doen't exists this function will create one.
func CreatePersistence(input *CreatePersistenceInput) Credentials {
	var user string
	var users []string
	var err error

	svc, err := createClient(
		input.AccessKey,
		input.SecretKey,
		input.SessionToken,
	)

	if err != nil {
		log.Fatalf("invalid creds error: %v", err)
	}

	if input.Username == "" {
		us, err := svc.listUsers()
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			}
		}
		users = us
		randIndex := rand.Intn(len(users))
		user = users[randIndex]
	} else {
		user = input.Username
		err = svc.createUser(user)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case iam.ErrCodeEntityAlreadyExistsException:
					fmt.Printf("user %s alredy exists", input.Username)
				default:
					fmt.Println(aerr.Error())
				}
			}
		}
	}

	err = svc.attachAdminPolicy(user)
	if err != nil {
		fmt.Println(err)
	}

	creds, err := svc.createAccessKey(user)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeLimitExceededException:
				fmt.Printf("error limit exceeded: %v", aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		}
	}

	return creds
}
