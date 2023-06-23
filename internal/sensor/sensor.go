package sensor

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/MichaelS11/go-dht"
)

// ConnectSensor initializes and opens DHT sensor connected to the GPIO.
func ConnectSensor(pinName, dhtModel string, fahrenheit bool) (*dht.DHT, error) {
	err := dht.HostInit()
	if err != nil {
		return nil, fmt.Errorf("sensor init error: %w", err)

	}

	unit := dht.Celsius
	if fahrenheit {
		unit = dht.Fahrenheit
	}

	dht, err := dht.NewDHT(pinName, unit, dhtModel)
	if err != nil {
		return nil, fmt.Errorf("sensor connection error: %w", err)
	}

	return dht, nil
}

// WatchSensor gets message from the sensor and send it to the `ch` channel.
//
// The function is blocking. To stop it, the context must be canceled.
func WatchSensor(ctx context.Context, dht *dht.DHT, ch chan<- Payload) {
	log.Info().Msg("watching sensor")
	for {
		timer := time.NewTimer(10 * time.Second)

		select {
		case <-timer.C:
			if humidity, temperature, err := dht.ReadRetry(3); err != nil {
				log.Warn().Err(err).Msg("sensor read error")
			} else {
				log.Debug().Msgf("sensor read: temperature='%f' humidity='%f'", temperature, humidity)
				ch <- Payload{Temperature: temperature, Humidity: humidity}
			}
		case <-ctx.Done():
			return
		}
	}
}
