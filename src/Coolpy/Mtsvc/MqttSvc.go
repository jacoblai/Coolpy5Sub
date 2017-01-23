package Mtsvc

import (
	"github.com/jacoblai/client"
	"github.com/jacoblai/transport"
	"github.com/jacoblai/broker"
	"strconv"
	"fmt"
)

type MqttSvc struct {
	Engine *broker.Engine
}

var Mport, Wsport int

func (m *MqttSvc) Host(mport, wsport int) {
	Mport = mport
	Wsport = wsport
	TcpServer, err := transport.Launch("tcp://:" + strconv.Itoa(Mport))
	if err != nil {
		panic(err)
	}
	WsServer, err := transport.Launch("ws://:" + strconv.Itoa(Wsport))
	if err != nil {
		panic(err)
	}
	backend := broker.NewMemoryBackend()
	m.Engine = broker.NewEngineWithBackend(backend)
	m.Engine.Accept(TcpServer)
	m.Engine.Accept(WsServer)
}

func Public(k string, payload []byte) {
	client := client.New()
	defer client.Close()
	cf, err := client.Connect("tcp://127.0.0.1:" + strconv.Itoa(Mport), nil)
	if err != nil {
		fmt.Println(err)
	}
	cf.Wait()
	pf, err := client.Publish(k, payload, 0, false)
	if err != nil {
		fmt.Println(err)
	}
	pf.Wait()
	err = client.Disconnect()
	if err != nil {
		fmt.Println(err)
	}
}