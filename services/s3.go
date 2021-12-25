package services

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func ConvertPayloadToS3Event(payload []byte) events.S3Event {
	var event events.S3Event
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err)
	}
	return event
}

//goland:noinspection GoUnusedParameter
func HandleS3Event(event events.S3Event, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
