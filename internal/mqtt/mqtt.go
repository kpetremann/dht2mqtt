package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const clientID = "dht2mqtt"
const mqttKeepAlive = 60
const mqttTimeout = 5

type Publisher struct {
	client    mqtt.Client
	mqttUrl   string
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

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info().Msg("connected to MQTT")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Error().Err(err).Msg("connecttion to MQTT lost")
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
	opts.SetOnConnectHandler(connectHandler)
	opts.SetConnectionLostHandler(connectLostHandler)

	p.client = mqtt.NewClient(opts)

	if t := p.client.Connect(); t.WaitTimeout(10*time.Second) && t.Error() != nil {
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

	log.Debug().Msgf("publishing new MQTT message: '%s %s'", p.topicRoot+sensorName, data)
	t := p.client.Publish(p.topicRoot+sensorName, 1, false, data)
	ack := t.WaitTimeout(10 * time.Second)

	if !ack {
		return errors.New("MQTT server did not confirm receiving the message")
	}
	if t.Error() != nil {
		return fmt.Errorf("MQTT server publish error: %w", t.Error())
	}

	log.Debug().Msgf("published: '%s %s'", p.topicRoot+sensorName, data)

	return nil
}
