package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kpetremann/dht2mqtt/internal/mqtt"
	"github.com/kpetremann/dht2mqtt/internal/sensor"
	"github.com/rs/zerolog/log"
)

func watchAndPublish(sensorName string, publisher mqtt.Publisher, ch <-chan sensor.Payload) {
	var lastPayload sensor.Payload
	lastChange := time.Now()
	for sensorPayload := range ch {
		if !sensorPayload.EqualTo(lastPayload) || time.Since(lastChange) > 5*time.Minute {
			lastPayload = sensorPayload
			err := publisher.Publish(sensorName, sensorPayload)
			if err != nil {
				log.Error().Err(err).Msg("failed to publish message to MQTT")
			}
			lastChange = time.Now()
		}
	}
	publisher.Disconnect()
}

func main() {
	cfg := readConfig()
	configureLogging(cfg.LogLevel)

	// Connect to the DHT sensor
	log.Info().Msgf("connecting to sensor on %s", cfg.GPIOPinName)
	dht, err := sensor.ConnectSensor(cfg.GPIOPinName, cfg.DHTmodel, cfg.Fahrenheit)
	if err != nil {
		log.Fatal().Err(err).Msg("new sensor error")
	}
	log.Info().Msg("connected to sensor")

	// Connect to MQTT server
	publisher := mqtt.NewPublisher(cfg.MQTTUrl, cfg.MQTTTopicRoot, cfg.MQTTUsername, cfg.MQTTPassword)
	if err := publisher.Connect(); err != nil {
		log.Fatal().Err(err).Msg("publisher init error")
	}

	// Watch for metrics
	ch := make(chan sensor.Payload, 10)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Info().Msg("waiting for metrics to send")
		watchAndPublish(cfg.SensorName, publisher, ch)
		log.Warn().Msg("publisher stopped")
		wg.Done()
	}()

	// Send metrics to MQTT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Info().Msg("start watching sensor")
	sensor.WatchSensor(ctx, dht, ch)
	log.Warn().Msg("stop watching sensor")

	// Stopping
	close(ch)
	wg.Wait()
}
