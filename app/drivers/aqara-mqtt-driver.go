package drivers

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type AqaraShutter interface {
	SetPosition(position float64) bool
}

type aqaraShutter struct {
	mqttClient mqtt.Client
	mqttTopic  string
}

func CreateAqaraShutter(mqttClient mqtt.Client, mqttTopic string) AqaraShutter {
	return &aqaraShutter{
		mqttClient: mqttClient,
		mqttTopic:  mqttTopic,
	}
}

type aqaraPositionRequest struct {
	Position int `json:"position"`
}

func (as *aqaraShutter) SetPosition(position float64) bool {
	body, _ := json.Marshal(aqaraPositionRequest{
		Position: int(position * 100),
	})

	token := as.mqttClient.Publish(as.mqttTopic, 0, false, body)

	if token.Wait() && token.Error() != nil {
		return false
	}

	return true
}
