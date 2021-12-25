package services

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ohoareau/gola/common"
	"log"
)

func ConvertPayloadToKinesisEvent(payload []byte) events.KinesisEvent {
	var event events.KinesisEvent
	err := json.Unmarshal(payload, &event)
	if nil != err {
		fmt.Println(err)
	}
	return event
}

//goland:noinspection GoUnusedParameter
func HandleKinesisEvent(event events.KinesisEvent, ctx interface{}, options common.Options) (interface{}, error) {
	for _, r := range event.Records {
		log.Println(r)
	}
	return nil, nil
}
