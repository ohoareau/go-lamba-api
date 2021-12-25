package services

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func ConvertPayloadToSnsEvent(payload []byte) events.SNSEvent {
	var event events.SNSEvent
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err.Error())
	}
	return event
}

//goland:noinspection GoUnusedParameter
func HandleSnsEvent(event events.SNSEvent, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
