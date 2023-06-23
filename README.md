# DHT2MQTT

DHT2MQTT is a yet another tool to send DHT11/DHT22 metrics to MQTT.

It is designed to be compatible with [mqtt-exporter](https://github.com/kpetremann/mqtt-exporter).

Only tested on Raspberry Pi 3B.

It leverages [MichaelS11/go-dht](https://github.com/MichaelS11/go-dht).

## Quickstart

### Build from source

```CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build ./cmd/dht2mqtt```

### Usage

```
./dht2mqtt -h
Usage of ./dht2mqtt:
  -gpio-pin string
    	GPIO PIN Name on which the sensor is connected (default "GPIO4")
  -mqtt-topic-root string
    	MQTT url, example: dh2mqtt/ (default "dh2mqtt/")
  -mqtt-url string
    	MQTT url, example: tcp://127.0.0.1:1883
  -mqtt-username string
    	username to connect to MQTT. The password must be set as in 'DHT2MQTT_PASSWORD' varenv
  -sensor-name string
    	sensor name (default "sensor")
```
