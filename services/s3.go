package services

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func IsS3Event(event interface{}) bool {
	switch event.(type) {
	case events.S3Event:
		return true
	default:
		return false
	}
}

//goland:noinspection GoUnusedParameter
func HandleS3Event(event events.S3Event, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
