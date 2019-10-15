package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func sanitize(value *string) {
	*value = strings.TrimPrefix(*value, "'")
	*value = strings.TrimSuffix(*value, "'")
}

func main() {
	app := cli.NewApp()
	app.Name = "cupt"
	app.Version = "0.2.1"
	app.Usage = "a Cognito User Pool tool.  Wrap values with special characters in single quotes."

	var configurationPath string
	var userPoolID string
	var clientID string
	var email string
	var password string
	var file string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "credentials, c",
			Usage:       "path to AWS configuration JSON file",
			Destination: &configurationPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a Cognito User Pool User",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "poolid, p",
					Usage:       "the Cognito User Pool Id to use",
					Destination: &userPoolID,
				},
				cli.StringFlag{
					Name:        "clientid, C",
					Usage:       "the Cognito User Pool Client Id to use",
					Destination: &clientID,
				},
				cli.StringFlag{
					Name:        "email, e",
					Usage:       "the email of the user to add",
					Destination: &email,
				},
				cli.StringFlag{
					Name:        "password, P",
					Usage:       "the password of the user to add",
					Destination: &password,
				},
			},
			Before: func(c *cli.Context) error {
				sanitize(&configurationPath)
				sanitize(&userPoolID)
				sanitize(&clientID)
				sanitize(&email)
				sanitize(&password)

				return nil
			},
			Action: func(c *cli.Context) error {
				svc := GetCognitoService(configurationPath)
				resp := AddUser(userPoolID, clientID, email, password, svc)
				log.Printf("%v", resp)
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "Get all Cognito User Pool Users",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "poolid, p",
					Usage:       "the Cognito User Pool Id to use",
					Destination: &userPoolID,
				},
			},
			Before: func(c *cli.Context) error {
				sanitize(&configurationPath)
				sanitize(&userPoolID)

				return nil
			},
			Action: func(c *cli.Context) error {
				svc := GetCognitoService(configurationPath)
				users := ListUsers(userPoolID, svc)
				log.Printf("%v", users)
				return nil
			},
		},
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "Log in as a Cognito User Pool User",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "poolid, p",
					Usage:       "the Cognito User Pool Id to use",
					Destination: &userPoolID,
				},
				cli.StringFlag{
					Name:        "clientid, C",
					Usage:       "the Cognito User Pool Client Id to use",
					Destination: &clientID,
				},
				cli.StringFlag{
					Name:        "email, e",
					Usage:       "the email of the user to add",
					Destination: &email,
				},
				cli.StringFlag{
					Name:        "password, P",
					Usage:       "the password of the user to add",
					Destination: &password,
				},
			},
			Before: func(c *cli.Context) error {
				sanitize(&configurationPath)
				sanitize(&userPoolID)
				sanitize(&clientID)
				sanitize(&email)
				sanitize(&password)

				return nil
			},
			Action: func(c *cli.Context) error {
				svc := GetCognitoService(configurationPath)
				resp := Login(userPoolID, clientID, email, password, svc)
				log.Printf("%v", resp)
				return nil
			},
		},
		{
			Name:    "backup",
			Aliases: []string{"b"},
			Usage:   "Serialize all Cognito User Pool Users to a file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "poolid, p",
					Usage:       "the Cognito User Pool Id to use",
					Destination: &userPoolID,
				},
				cli.StringFlag{
					Name:        "file, f",
					Usage:       "the name of a file to serialize users to",
					Destination: &file,
				},
			},
			Before: func(c *cli.Context) error {
				sanitize(&configurationPath)
				sanitize(&userPoolID)
				sanitize(&file)

				return nil
			},
			Action: func(c *cli.Context) error {
				log.Printf("Beginning backup process.  This may take a while (about 5 seconds per 100 users).")
				svc := GetCognitoService(configurationPath)
				users := ListUsers(userPoolID, svc)
				len := WriteUsers(file, users)
				log.Printf("Serialized %v users to %v.", len, file)
				return nil
			},
		},
		{
			Name:    "restore",
			Aliases: []string{"r"},
			Usage:   "Restore a Cognito User Pool's users from a serialized JSON file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "poolid, p",
					Usage:       "the Cognito User Pool Id to use",
					Destination: &userPoolID,
				},
				cli.StringFlag{
					Name:        "file, f",
					Usage:       "the name of a file to serialize users to",
					Destination: &file,
				},
			},
			Before: func(c *cli.Context) error {
				sanitize(&configurationPath)
				sanitize(&userPoolID)
				sanitize(&file)

				return nil
			},
			Action: func(c *cli.Context) error {
				log.Printf("Beginning restoration process.  This may take quite a long time (about 50 seconds per 100 users).")
				svc := GetCognitoService(configurationPath)
				users := ReadUsers(file)
				successes := 0
				for _, user := range users {
					if RestoreUser(userPoolID, user, svc) {
						successes++
					}
				}
				log.Printf("Restoration completed.  Restored %v of %v users successfully.", successes, len(users))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
