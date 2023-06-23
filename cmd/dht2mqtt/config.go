package main

import (
	"flag"
	"log"
	"os"
)

type config struct {
	SensorName    string
	GPIOPinName   string
	DHTmodel      string
	MQTTUrl       string
	MQTTTopicRoot string
	MQTTUsername  string
	MQTTPassword  string
	Fahrenheit    bool
}

func readConfig() config {
	sensorName := flag.String("sensor-name", "sensor", "sensor name")
	gpioPinName := flag.String("gpio-pin", "GPIO4", "GPIO PIN Name on which the sensor is connected")
	dhtModel := flag.String("dht-model", "DHT22", "DHT sensor model: DHT11 or DHT22")
	fahrenheit := flag.Bool("fahrenheit", false, "Temperature unit. Fahrenheit if set, default is Celcius")

	mqttUrl := flag.String("mqtt-url", "", "MQTT url, example: tcp://127.0.0.1:1883")
	mqttTopicRoot := flag.String("mqtt-topic-root", "dh2mqtt/", "MQTT url, example: dh2mqtt/")

	mqttUsername := flag.String("mqtt-username", "", "username to connect to MQTT. The password must be set as in 'DHT2MQTT_PASSWORD' varenv")
	mqttPassword := os.Getenv("DHT2MQTT_PASSWORD")

	flag.Parse()

	if *mqttUrl == "" {
		log.Fatalln("mqtt-url undefined")
	}

	if *mqttTopicRoot == "" {
		log.Fatalln("topicBase undefined")
	}

	return config{
		SensorName:    *sensorName,
		GPIOPinName:   *gpioPinName,
		DHTmodel:      *dhtModel,
		Fahrenheit:    *fahrenheit,
		MQTTUrl:       *mqttUrl,
		MQTTTopicRoot: *mqttTopicRoot,
		MQTTUsername:  *mqttUsername,
		MQTTPassword:  mqttPassword,
	}
}
