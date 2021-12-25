package services

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func IsKinesisEvent(event interface{}) bool {
	switch event.(type) {
	case events.KinesisEvent:
		return true
	default:
		return false
	}
}

//goland:noinspection GoUnusedParameter
func HandleKinesisEvent(event events.KinesisEvent, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
