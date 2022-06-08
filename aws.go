package iampersist

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	AdministratorAccessARN = "arn:aws:iam::aws:policy/AdministratorAccess"
)

type Client struct {
	*iam.IAM
}

type Credentials struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// createClient - create iam client or error
func createClient(accesskey, secretkey, sessiontoken string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accesskey, secretkey, sessiontoken),
	})

	if err != nil {
		return nil, err
	}

	if _, err := sess.Config.Credentials.Get(); err != nil {
		return nil, err
	}

	svc := iam.New(sess)

	return &Client{svc}, nil
}

// listUsers - return iam users or error
func (c *Client) listUsers() ([]string, error) {

	users, err := c.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return nil, err
	}

	var usernames []string

	for _, user := range users.Users {
		usernames = append(usernames, *user.UserName)
	}

	return usernames, nil
}

// createAccessKey - create aws creds and return struct or error
func (c *Client) createAccessKey(arn string) (Credentials, error) {

	req, resp := c.CreateAccessKeyRequest(&iam.CreateAccessKeyInput{UserName: aws.String(arn)})

	err := req.Send()
	if err == nil {
		return Credentials{
			AccessKey: *resp.AccessKey.AccessKeyId,
			SecretKey: *resp.AccessKey.SecretAccessKey,
		}, nil
	}

	return Credentials{}, err

}

// attachAdminPolicy - attch Admistrator Policy or return error
func (c *Client) attachAdminPolicy(username string) error {

	result, err := c.AttachUserPolicy(&iam.AttachUserPolicyInput{
		UserName:  aws.String(username),
		PolicyArn: aws.String(AdministratorAccessARN),
	})

	if err != nil {
		return err
	}

	_ = result

	return nil
}

// createUser - create iam user or return error
func (c *Client) createUser(username string) error {

	result, err := c.CreateUser(&iam.CreateUserInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return err
	}

	if result.User.Arn != nil {
		return nil
	}

	return nil
}
