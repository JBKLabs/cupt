package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// WriteUsers serializes a slice of users to a JSON file.
func WriteUsers(file string, users []*cognitoidentityprovider.UserType) int {
	json, jsonErr := json.Marshal(users)
	if jsonErr != nil {
		log.Fatal("Failed to encode the users to JSON: ", jsonErr.Error())
	}

	err := ioutil.WriteFile(file, json, 0644)
	if err != nil {
		log.Fatal("Failed to write the users to a file: ", err.Error())
	}

	return len(users)
}

// ReadUsers deserializes an array of users from a JSON file.
func ReadUsers(file string) []*cognitoidentityprovider.UserType {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Failed to read the JSON file: ", err.Error())
		os.Exit(1)
	}

	var users []*cognitoidentityprovider.UserType
	json.Unmarshal(raw, &users)

	// We must filter out the "sub" property that is managed by Cognito.
	for userIndex, user := range users {
		for attrIndex, attr := range user.Attributes {
			if *attr.Name == "sub" {
				user.Attributes = append(user.Attributes[:attrIndex], user.Attributes[attrIndex+1:]...)
				break
			}
		}

		users[userIndex] = user
	}

	return users
}
