package services

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func ConvertPayloadToDynamoDBEvent(payload []byte) events.DynamoDBEvent {
	var event events.DynamoDBEvent
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err.Error())
	}
	return event
}

//goland:noinspection GoUnusedParameter
func HandleDynamoDBEvent(event events.DynamoDBEvent, ctx interface{}, options *common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
