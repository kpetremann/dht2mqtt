# DHT2MQTT

DHT2MQTT is yet another tool to send DHT11/DHT22 metrics to MQTT.

It is designed to be compatible with [mqtt-exporter](https://github.com/kpetremann/mqtt-exporter).

Only tested on Raspberry Pi 3B.

It leverages [MichaelS11/go-dht](https://github.com/MichaelS11/go-dht).

## Quickstart

### Build from source

```CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build ./cmd/dht2mqtt```

### Usage

```
Usage of ./dht2mqtt:
  -dht-model string
        DHT sensor model: DHT11 or DHT22 (default "DHT22")
  -fahrenheit
        Temperature unit. Fahrenheit if set, default is Celcius
  -gpio-pin string
        GPIO PIN Name on which the sensor is connected (default "GPIO4")
  -log-level string
        log level (debug, info, warn, error, fatal, panic, disabled) (default "info")
  -mqtt-topic-root string
        MQTT url, example: dh2mqtt/ (default "dh2mqtt/")
  -mqtt-url string
        MQTT url, example: tcp://127.0.0.1:1883
  -mqtt-username string
        username to connect to MQTT. The password must be set as in 'DHT2MQTT_PASSWORD' varenv
  -sensor-name string
        sensor name (default "sensor")
```

Example:
```
~ $ DHT2MQTT_PASSWORD="awesomepassword" ./dht2mqtt -mqtt-url tcp://10.0.0.1:1883 -sensor-name garden -mqtt-username dht
```

### Systemd service example

* Put the binary somewhere like in `/usr/local/bin/`.
* Create the file `/etc/systemd/system/dht2mqtt` with the following content:

```
[Unit]
Description = dht2mqtt
Wants = network-online.target
After = network-online.target

[Install]
WantedBy = multi-user.target

[Service]
Type = simple
Environment = "DHT2MQTT_PASSWORD=<awesomepassword>"
ExecStart = /usr/local/bin/dht2mqtt -mqtt-url tcp://10.0.0.:1883 -sensor-name <sensor-name> -mqtt-username <user>
```

**Do not forget to adapt the variables.**

* Then simply run it and enable on boot:
```
sudo systemctl daemon-reload

systemctl enable --now dht2mqtt
# or:
sudo systemctl enable dht2mqtt
sudo systemctl start dht2mqtt
```