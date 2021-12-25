package services

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func IsSnsEvent(event interface{}) bool {
	switch event.(type) {
	case events.SNSEvent:
		return true
	default:
		return false
	}
}

//goland:noinspection GoUnusedParameter
func HandleSnsEvent(event events.SNSEvent, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
