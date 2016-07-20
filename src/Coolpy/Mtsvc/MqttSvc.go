package Mtsvc

import (
	"github.com/surgemq/surgemq/service"
	"log"
	"strconv"
	"github.com/surgemq/message"
)

var svc *service.Server

func Host(mport int) {
	// Create a mqtt server
	//auth.Register("coolpy", &Manager{})
	mqttsvr := &service.Server{
		KeepAlive:        300, // seconds
		ConnectTimeout:   2, // seconds
		SessionsProvider: "mem", // keeps sessions in memory
		Authenticator:    "mockSuccess", // always succeed
		TopicsProvider:   "mem", // keeps topic subscriptions in memory
	}
	svc = mqttsvr
	go func() {
		// Listen and serve connections at mport
		if err := mqttsvr.ListenAndServe("tcp://:" + strconv.Itoa(mport)); err != nil {
			log.Fatal(err)
		}
	}()
}

func Close() {
	svc.Close()
}

func Public(k string,payload []byte) {
	// Creates a new PUBLISH message with the appropriate contents for publishing
	pubmsg := message.NewPublishMessage()
	pubmsg.SetPacketId(1)
	pubmsg.SetTopic([]byte(k))
	pubmsg.SetPayload(payload)
	pubmsg.SetQoS(0)
	// Publishes to the server by sending the message
	svc.Publish(pubmsg, nil)
}