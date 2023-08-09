package main

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// MQTT broker address and port
	brokerAddress := "clientId-WgKlWeF4if:8884"

	// Create a new MQTT client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerAddress)
	client := mqtt.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Create the JSON message
	message := map[string]interface{}{
		"messageID": "interfaces",
		"data": map[string]string{
			"interface": "Wi-Fi",
			"mac":       "B0-A4-60-EB-2E-40",
		},
	}

	// Convert JSON to string
	messageString, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	// Publish the message to a topic
	topic := "Communication" // Replace with your desired topic
	token := client.Publish(topic, 0, false, messageString)
	token.Wait()

	fmt.Println("Message sent:", string(messageString))

	// Disconnect from the broker
	client.Disconnect(250)
}
