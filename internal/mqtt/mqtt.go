package mqtt

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/zerolog/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const clientID = "dht2mqtt"
const mqttKeepAlive = 60
const mqttTimeout = 5

type Publisher struct {
	client    mqtt.Client
	mqttUrl   string // example: "tcp://10.2.0.166:1883"
	username  string
	password  string
	topicRoot string
}

// NewPublisher instantiates a MQTT publisher.
//
// It connects to `mqttUrl` with `username` and `password` if specified.
// TopicRoot will be used as prefix of topic for all MQTT messages.
func NewPublisher(mqttUrl, topicRoot, username, password string) Publisher {
	return Publisher{
		mqttUrl:   mqttUrl,
		topicRoot: topicRoot,
		username:  username,
		password:  password,
	}
}

// Connect connects to the MQTT server;
//
// The Publisher must have been instantiated first using NewPublisher.
func (p *Publisher) Connect() error {
	opts := mqtt.NewClientOptions().AddBroker(p.mqttUrl).SetClientID(clientID)
	opts.SetKeepAlive(mqttKeepAlive * time.Second)
	opts.SetPingTimeout(mqttTimeout * time.Second)
	if p.username != "" && p.password != "" {
		opts.SetUsername(p.username)
		opts.SetPassword(p.password)
	}
	opts.SetAutoReconnect(true)

	p.client = mqtt.NewClient(opts)

	if t := p.client.Connect(); t.Wait() && t.Error() != nil {
		log.Fatal().Err(t.Error()).Send()
	}

	return nil
}

func (p *Publisher) Disconnect() {
	p.client.Disconnect(250)
	time.Sleep(1 * time.Second)
}

// Publish sends the message to MQTT server.
//
// The `sensorName` will be in the topic.
// The payload must be JSON serializable.
func (p *Publisher) Publish(sensorName string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("publishing new MQTT message: '%s %s'", p.topicRoot+sensorName, data)
	t := p.client.Publish(p.topicRoot+sensorName, 0, false, data)
	if ok := t.Wait(); !ok {
		return errors.New("MQTT server did not confirm receiving the message")
	}

	return nil
}
