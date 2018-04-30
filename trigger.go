package mytrigger

import (
	"context"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-twitter-listner")

// TwitterTrigger is simple MQTT trigger
type TwitterTrigger struct {
	metadata       *trigger.Metadata
	config         *trigger.Config
	handlers       []*trigger.Handler
	topicToHandler map[string]*trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &TwitterFactory{metadata: md}
}

// TwitterFactory MQTT Trigger factory
type TwitterFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *TwitterFactory) New(config *trigger.Config) trigger.Trigger {
	return &TwitterTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *TwitterTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Initialize implements trigger.Initializable.Initialize
func (t *TwitterTrigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	return nil
}

// Start implements trigger.Trigger.Start
func (t *TwitterTrigger) Start() error {

	t.topicToHandler = make(map[string]*trigger.Handler)

	apiKey := strings.TrimSpace(t.config.GetSetting("apiKey"))
	apiSecret := strings.TrimSpace(t.config.GetSetting("apiSecret"))
	consumerKey := strings.TrimSpace(t.config.GetSetting("consumerKey"))
	consumerSecret := strings.TrimSpace(t.config.GetSetting("consumerSecret"))
	searchString := strings.TrimSpace(t.handlers[0].GetStringSetting("searchString"))
	if len(apiKey) == 0 || len(apiSecret) == 0 || len(consumerKey) == 0 || len(consumerSecret) == 0 || len(searchString) == 0 {
		log.Info("Please check the input parameters")
		panic("Error Occured Due To Input Params")
	} else {

		log.Info("")

		for _, handler := range t.handlers {

			topic := handler.GetStringSetting("searchString")

			//	t.RunHandler(handler, "test1")

			anaconda.SetConsumerKey(apiKey)
			anaconda.SetConsumerSecret(apiSecret)
			api := anaconda.NewTwitterApi(consumerKey, consumerSecret)

			stream := api.PublicStreamFilter(url.Values{
				"track": []string{topic},
			})

			defer stream.Stop()

			for v := range stream.C {
				twt, _ := v.(anaconda.Tweet)
				log.Info("Received Tweet", twt.IdStr, twt.User.Name, twt.Text)
				t.RunHandler(handler, twt.IdStr, twt.User.Name, twt.Text)
			}

			// ticker := time.NewTicker(60 * time.Second)
			// go func() {
			// 	for te := range ticker.C {
			// 		s := time.Now().String()
			// 		fmt.Println(te.Hour())
			// 		t.RunHandler(handler, s)
			// 	}

			//	}()

			log.Debugf("topic: [%s]", topic)

		}

	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *TwitterTrigger) Stop() error {
	//unsubscribe from topic

	return nil
}

// RunHandler runs the handler and associated action
func (t *TwitterTrigger) RunHandler(handler *trigger.Handler, id string, name string, message string) {

	trgData := make(map[string]interface{})
	trgData["tweetId"] = id
	trgData["screenName"] = name
	trgData["message"] = message

	results, err := handler.Handle(context.Background(), trgData)

	if err != nil {
		log.Error("Error starting action: ", err.Error())
	}
	log.Debugf("Ran Handler: [%s]", results)
	log.Debugf("Ran Handler: [%s]", handler)

}
