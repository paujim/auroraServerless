package main

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

type DataServiceAPI interface {
	ExecuteStatement(input *rdsdataservice.ExecuteStatementInput) (*rdsdataservice.ExecuteStatementOutput, error)
}

type SqlClient struct {
	client    DataServiceAPI
	auroraArn *string
	secretArn *string
}

func (c *SqlClient) InsertProfile(fullName, email string, phoneNumbers []string) (*int64, error) {
	log.Printf("Insert data to DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: c.auroraArn,
		SecretArn:   c.secretArn,
		Sql:         aws.String("INSERT INTO TestDB.Profiles (FullName, Email, Phones) VALUES (:name, :email, :phones);"),
		Parameters: []*rdsdataservice.SqlParameter{
			{
				Name: aws.String("name"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(fullName),
				},
			},
			{
				Name: aws.String("email"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(email),
				},
			},
			{
				Name: aws.String("phones"),
				Value: &rdsdataservice.Field{
					StringValue: aws.String(strings.Join(phoneNumbers, ";")),
				},
			},
		},
	}
	resp, err := c.client.ExecuteStatement(params)
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return nil, err
	}
	log.Printf("%s\n", resp.GoString())
	return resp.GeneratedFields[0].LongValue, nil
}

func (h *SqlClient) GetProfiles() ([]Profile, error) {
	log.Printf("Get data from DB\n")

	params := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: h.auroraArn,
		SecretArn:   h.secretArn,
		Sql:         aws.String("SELECT * FROM TestDB.Profiles"),
	}
	resp, err := h.client.ExecuteStatement(params)
	if err != nil {
		log.Printf("Error fetching profiles: %s", err)
		return nil, err
	}

	profiles := []Profile{}
	for _, record := range resp.Records {
		profiles = append(profiles, Profile{
			ID:           *record[0].LongValue,
			FullName:     *record[NAME].StringValue,
			Email:        *record[EMAIL].StringValue,
			PhoneNumbers: []string{*record[PHONE].StringValue},
		})
	}
	return profiles, nil
}
