package mytrigger

import (
	"context"
	"fmt"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-mqtt")

// MqttTrigger is simple MQTT trigger
type MqttTrigger struct {
	metadata       *trigger.Metadata
	config         *trigger.Config
	handlers       []*trigger.Handler
	topicToHandler map[string]*trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MQTTFactory{metadata: md}
}

// MQTTFactory MQTT Trigger factory
type MQTTFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *MQTTFactory) New(config *trigger.Config) trigger.Trigger {
	return &MqttTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *MqttTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Initialize implements trigger.Initializable.Initialize
func (t *MqttTrigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	return nil
}

// Start implements trigger.Trigger.Start
func (t *MqttTrigger) Start() error {

	t.topicToHandler = make(map[string]*trigger.Handler)

	log.Info("Start called")

	for _, handler := range t.handlers {

		topic := handler.GetStringSetting("topic")

		//	t.RunHandler(handler, "test1")

		ticker := time.NewTicker(60 * time.Second)
		go func() {
			for te := range ticker.C {
				s := time.Now().String()
				fmt.Println(te.Hour())
				t.RunHandler(handler, s)
			}

		}()

		log.Debugf("topic: [%s]", topic)

	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *MqttTrigger) Stop() error {
	//unsubscribe from topic

	return nil
}

// RunHandler runs the handler and associated action
func (t *MqttTrigger) RunHandler(handler *trigger.Handler, payload string) {

	trgData := make(map[string]interface{})
	trgData["message"] = payload

	results, err := handler.Handle(context.Background(), trgData)

	if err != nil {
		log.Error("Error starting action: ", err.Error())
	}
	log.Debugf("Ran Handler: [%s]", results)
	log.Debugf("Ran Handler: [%s]", handler)

}
