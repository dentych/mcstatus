package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var previousStatus = true

func main() {
	mcServerUri := os.Getenv("MC_STATUS_SERVER_URL")
	if mcServerUri == "" {
		log.Fatalf("MC_SERVER_URI is not set but is required.")
	}

	discordWebhook := os.Getenv("MC_STATUS_DISCORD_WEBHOOK")
	if discordWebhook == "" {
		log.Fatalf("MC_STATUS_DISCORD_WEBHOOK is not set but is required.")
	}

	for {
		time.Sleep(10 * time.Second)
		isUp, err := runCheck(mcServerUri)
		if err != nil {
			log.Printf("Failed to run check: %s\n", err)
		}
		if isUp && !previousStatus {
			err = celebrateOnDiscord(discordWebhook)
			if err != nil {
				log.Printf("Failed to celebrate on Discord: %s\n", err)
				continue
			} else {
				previousStatus = true
			}
		} else if !isUp && previousStatus {
			err = complainOnDiscord(discordWebhook)
			if err != nil {
				log.Printf("Failed to complain on Discord: %s\n", err)
				continue
			} else {
				previousStatus = false
			}
		}
	}

	fmt.Println("Works")
}

func celebrateOnDiscord(webhookUrl string) error {
	log.Printf("Discord: Server is UP!")

	marshalled, err := json.Marshal(map[string]string{
		"content": ":white_check_mark: Minecraft server is UP! :clap:",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal discord message: %w", err)
	}
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewReader(marshalled))
	if err != nil {
		return fmt.Errorf("failed to send post response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("response code was unsuccessful, statusCode=%d", resp.StatusCode)
	}

	return nil
}

func complainOnDiscord(webhookUrl string) error {
	log.Printf("Discord: Server is UP!")

	marshalled, err := json.Marshal(map[string]string{
		"content": ":no_entry_sign: Minecraft server is DOWN! :rotating_light:",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal discord message: %w", err)
	}
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewReader(marshalled))
	if err != nil {
		return fmt.Errorf("failed to send post response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("response code was unsuccessful, statusCode=%d", resp.StatusCode)
	}

	return nil
}

func runCheck(mcServerUri string) (bool, error) {
	addr, err := net.ResolveTCPAddr("tcp", mcServerUri)
	if err != nil {
		log.Printf("Failed to resolve TCP address: %s\n", err)
		return false, err
	}

	con, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Printf("Failed to dial TCP connection: %s\n", err)
		return false, nil
	}

	err = con.Close()
	if err != nil {
		log.Printf("Failed to close connection: %s\n", err)
	}

	return true, nil
}
