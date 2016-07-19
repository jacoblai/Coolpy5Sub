package MqttClient

import (
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
	"strconv"
)

var mqtt *service.Client

func Connect(port int) {
	// Instantiates a new Client
	mqtt = &service.Client{}
	// Creates a new MQTT CONNECT message and sets the proper parameters
	msg := message.NewConnectMessage()
	msg.SetWillQos(1)
	msg.SetVersion(4)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte("CoolpyService"))
	msg.SetKeepAlive(10)
	msg.SetWillTopic([]byte("will"))
	msg.SetWillMessage([]byte("send me home"))
	msg.SetUsername([]byte("coolpy"))
	msg.SetPassword([]byte("icoolpy.com"))

	// Connects to the remote server at 127.0.0.1 port 1883
	mqtt.Connect("tcp://127.0.0.1:" + strconv.Itoa(port), msg)

	// Disconnects from the server
	//mqtt.Disconnect()
}

func Public(topic string,payload string) {
	// Creates a new PUBLISH message with the appropriate contents for publishing
	pubmsg := message.NewPublishMessage()
	pubmsg.SetPacketId(3)
	pubmsg.SetTopic([]byte(topic))
	pubmsg.SetPayload([]byte(payload))
	pubmsg.SetQoS(0)
	// Publishes to the server by sending the message
	mqtt.Publish(pubmsg, nil)
}