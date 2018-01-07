package main

import (
	"log"

	. "github.com/mlabouardy/dialogflow-go-client"
	apiai "github.com/mlabouardy/dialogflow-go-client/models"
)

func GetResponse(input string, token string) apiai.Result {
	err, client := NewDialogFlowClient(apiai.Options{
		AccessToken: token,
	})
	if err != nil {
		log.Fatal(err)
	}

	query := apiai.Query{
		Query: input,
	}
	resp, err := client.QueryFindRequest(query)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Result
}
