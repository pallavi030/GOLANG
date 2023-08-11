package mqtt_project

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type InterfaceInfo struct {
	Interface string `json:"interface"`
	MAC       string `json:"mac"`
	IP        string `json:"ip"`
}

type MQTTConfig struct {
	URL      string `json:"mqtt_url"`
	Username string `json:"mqtt_username"`
	Password string `json:"mqtt_password"`
	Topic    string `json:"mqtt_topic"`
}

func LoadConfig(configPath string) (MQTTConfig, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return MQTTConfig{}, err
	}
	defer configFile.Close()

	var config MQTTConfig
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return MQTTConfig{}, err
	}

	return config, nil
}

func SendInterfacesToMQTT(configPath string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	opts := MQTT.NewClientOptions().AddBroker(config.URL)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetClientID("goClient")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250)

	interfaceInfo, err := getInterfaces()
	if err != nil {
		return err
	}

	message := MQTTMessage{
		MessageID:  "interfaces",
		Interfaces: interfaceInfo,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	token := client.Publish(config.Topic, 0, false, jsonMessage)
	token.Wait()

	fmt.Println("Message sent successfully.")
	return nil
}

func getInterfaces() ([]InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var interfaceInfo []InterfaceInfo

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}

			if ip.To4() != nil {
				interfaceInfo = append(interfaceInfo, InterfaceInfo{
					Interface: iface.Name,
					MAC:       iface.HardwareAddr.String(),
					IP:        ip.String(),
				})
			}
		}
	}

	return interfaceInfo, nil
}
