package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

// Struct to represent Discord webhook payload
type DiscordMessage struct {
	Content string `json:"content"`
}

func main() {
	// Function to run shell commands and return output
	runCommand := func(cmd string) string {
		parts := strings.Fields(cmd)
		head := parts[0]
		parts = parts[1:len(parts)]
		out, err := exec.Command(head, parts...).Output()
		if err != nil {
			return fmt.Sprintf("Error: %s", err)
		}
		return string(out)
	}

	// Execute commands to gather OS information
	osInfo1 := runCommand("uname -a")
	osInfo2 := runCommand("lsb_release -a")

	// Format the information
	message := fmt.Sprintf("OS Information for VM\n```\n%s```\n```\n%s\n```", osInfo1, osInfo2)

	// Send the message to Discord webhook
	sendToDiscord(message)
}

func sendToDiscord(message string) {
	// Discord webhook URL
	webhookURL := "{{webhookURL}}"

	// Create Discord message object
	discordMessage := DiscordMessage{
		Content: message,
	}

	// Convert Discord message to JSON
	jsonMessage, err := json.Marshal(discordMessage)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Send HTTP POST request to Discord webhook
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		fmt.Println("Error sending message to Discord:", err)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response status:", resp.StatusCode)
		return
	}

	fmt.Println("Message sent successfully to Discord!")
}
