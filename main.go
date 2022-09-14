package main

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"os"
	"path"
	"time"

	"curtain-controller/app/api"
	"curtain-controller/app/config"
	"curtain-controller/app/drivers"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	fmt.Println("Curtain controller starting")

	appConfig := config.LoadConfig(configPath("curtain-controller.yaml"))
	mqttClient := createMQTTClient(appConfig.MQTT)

	defer mqttClient.Disconnect(0)

	devices := drivers.DeviceList{
		AqaraShutters: drivers.AqaraShutters{},
		SlideCurtains: drivers.SlideCurtains{},
	}

	for n, c := range appConfig.Devices.AqaraShutters {
		devices.AqaraShutters[n] = drivers.CreateAqaraShutter(mqttClient, c.Topic)
	}

	for n, c := range appConfig.Devices.SlideCurtains {
		devices.SlideCurtains[n] = drivers.CreateSlideCurtain(c.IP, c.DeviceID)
	}

	fmt.Println("Curtain controller running")

	api.RunAPI(appConfig.API, devices)
}

func configPath(filePath string) string {
	exe, err := os.Executable()

	if err != nil {
		return filePath
	}

	return path.Join(path.Dir(exe), filePath)
}

func createMQTTClient(mqttConfig config.MQTTConfig) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", mqttConfig.Host, mqttConfig.Port)).
		SetClientID(mqttConfig.ClientId).
		SetUsername(mqttConfig.Username).
		SetPassword(mqttConfig.Password)

	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetAutoReconnect(true)

	// DEBUG LOG MESSAGES FOR MQTT CONNECTION STATUS
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("Connected to MQTT broker")
	}

	opts.OnConnectAttempt = func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		fmt.Printf("Attempting to connect to MQTT broker at %s\n", broker)

		return tlsCfg
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("Connection to MQTT broker lost: %s\n", err)
	}

	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Printf("Attempting to reconnect to MQTT broker\n")
	}

	// END DEBUG LOG MESSAGES FOR MQTT CONNECTION STATUS

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Unable to connect to MQTT broker: %s\n", token.Error())
	}

	return c
}
