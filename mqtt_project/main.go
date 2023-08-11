package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/pallavi030/GOLANG" // Replace with the correct module path
)


func main() {
	configPath := "mqtt_config.json"
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		os.Exit(1)
	}

	err = mqtt_project.SendInterfacesToMQTT(absConfigPath)
	if err != nil {
		fmt.Println("Error sending interfaces:", err)
		os.Exit(1)
	}
}
