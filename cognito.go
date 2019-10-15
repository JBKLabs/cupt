package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	uuid "github.com/nu7hatch/gouuid"
	pass "github.com/sethvargo/go-password/password"
)

// GetCognitoService constructs a cognitoidentityprovider service from an AWS configuration file.
func GetCognitoService(path string) cognitoidentityprovider.CognitoIdentityProvider {
	var config Config
	GetConfig(path, &config)
	creds := credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: creds,
	})
	if err != nil {
		log.Fatal("Failed to create a new AWS session: ", err.Error())
	}
	return *cognitoidentityprovider.New(sess)
}

// ListUsers lists the users in a Cognito User Pool.
func ListUsers(userPoolID string, svc cognitoidentityprovider.CognitoIdentityProvider) []*cognitoidentityprovider.UserType {
	token := ""
	finished := false
	var batch int64 = 5
	var users []*cognitoidentityprovider.UserType

	for !finished {
		args := &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(userPoolID),
			Limit:      aws.Int64(batch),
		}

		if token != "" {
			args.PaginationToken = aws.String(token)
		}

		resp, err := svc.ListUsers(args)

		if err != nil {
			log.Fatal("Failed to list users: ", err.Error())
		}
		users = append(users, resp.Users...)
		if resp.PaginationToken != nil && int64(len(resp.Users)) == batch {
			token = *resp.PaginationToken
		} else {
			finished = true
		}
	}

	return users
}

// RestoreUser restores a deserialized Cognito User Pool User.  They will be assigned a random password, but retain their username and attributes (excluding sub).
func RestoreUser(userPoolID string, user *cognitoidentityprovider.UserType, svc cognitoidentityprovider.CognitoIdentityProvider) bool {
	tempPass := pass.MustGenerate(64, 10, 10, false, false)
	_, err := svc.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        aws.String(userPoolID),
		TemporaryPassword: aws.String(tempPass),
		Username:          aws.String(*user.Username),
		UserAttributes:    user.Attributes,
	})

	if err != nil {
		log.Fatal("Failed to restore a user: ", *user.Username, " ", err.Error())
		return false
	}

	return true
}

// AddUser adds a user to the Cognito User Pool with a random GUID username and the supplied auto-verified email and password.
func AddUser(userPoolID string, clientID string, email string, password string, svc cognitoidentityprovider.CognitoIdentityProvider) *cognitoidentityprovider.AdminRespondToAuthChallengeOutput {
	tempPass := pass.MustGenerate(64, 10, 10, false, false)
	username, _ := uuid.NewV4()

	_, err := svc.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        aws.String(userPoolID),
		TemporaryPassword: aws.String(tempPass),
		Username:          aws.String(username.String()),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			&cognitoidentityprovider.AttributeType{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			&cognitoidentityprovider.AttributeType{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
		},
	})

	if err != nil {
		log.Fatal("failed to create the user: ", err.Error())
	}

	req, resp := svc.AdminInitiateAuthRequest(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   aws.String("ADMIN_NO_SRP_AUTH"),
		UserPoolId: aws.String(userPoolID),
		ClientId:   aws.String(clientID),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username.String()),
			"PASSWORD": aws.String(tempPass),
		},
	})

	err = req.Send()
	if err != nil {
		log.Fatal("failed to sign in as the newly created user: ", err)
	}

	if *resp.ChallengeName == "NEW_PASSWORD_REQUIRED" {
		session := resp.Session
		resp, err := svc.AdminRespondToAuthChallenge(&cognitoidentityprovider.AdminRespondToAuthChallengeInput{
			ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
			UserPoolId:    aws.String(userPoolID),
			ClientId:      aws.String(clientID),
			ChallengeResponses: map[string]*string{
				"USERNAME":     aws.String(username.String()),
				"NEW_PASSWORD": aws.String(password),
			},
			Session: session,
		})

		if err != nil {
			log.Fatal("failed to assign the newly created user's password: ", err)
		}

		return resp
	}

	return nil
}

// Login to a Cognito User Pool and get a set of tokens.
func Login(userPoolID string, clientID string, email string, password string, svc cognitoidentityprovider.CognitoIdentityProvider) *cognitoidentityprovider.AdminInitiateAuthOutput {
	req, resp := svc.AdminInitiateAuthRequest(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   aws.String("ADMIN_NO_SRP_AUTH"),
		UserPoolId: aws.String(userPoolID),
		ClientId:   aws.String(clientID),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(email),
			"PASSWORD": aws.String(password),
		},
	})

	err := req.Send()
	if err != nil {
		log.Fatal("failed to sign in: ", err)
	}

	return resp
}
