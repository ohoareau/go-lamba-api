package services

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func IsSqsEvent(event interface{}) bool {
	switch event.(type) {
	case events.SQSEvent:
		return true
	default:
		return false
	}
}

//goland:noinspection GoUnusedParameter
func HandleSqsEvent(event events.SQSEvent, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
