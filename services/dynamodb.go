package services

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func IsDynamoDBEvent(event interface{}) bool {
	switch event.(type) {
	case events.DynamoDBEvent:
		return true
	default:
		return false
	}
}

//goland:noinspection GoUnusedParameter
func HandleDynamoDBEvent(event events.DynamoDBEvent, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
