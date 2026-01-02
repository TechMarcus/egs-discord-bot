package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func SendGames(file string) error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	webhookUrl := os.Getenv("DISCORD_WEBHOOK_URL")
	fmt.Println("Sending to webhook:", webhookUrl)
	jsonfile, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	freeGames, err := FormatFreeGames(jsonfile)
	if err != nil {
		return err
	}
	payload, err := FreeGamesPayload(webhookUrl, freeGames)
	if err != nil {
		return err
	}
	fmt.Println(string(payload))
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	return nil
}

func FormatFreeGames(file []byte) ([]GameInfo, error) {
	var FreeGames []GameInfo
	json.Unmarshal(file, &FreeGames)
	return FreeGames, nil
}

type DiscordEmbed struct {
	Title string `json:"title"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
}

func FreeGamesPayload(webhookUrl string, games []GameInfo) ([]byte, error) {
	var payload = make(map[string][]DiscordEmbed)

	for _, game := range games {
		var e DiscordEmbed
		if game.Url != "" {
			e = DiscordEmbed{
				Title: game.Name + "\nLink: " + game.Url + "\nStart Date: " + game.StartDate + "\nEnd Date: " + game.EndDate,
			}
		} else {
			e = DiscordEmbed{
				Title: game.Name + "\nStart Date: " + game.StartDate + "\nEnd Date: " + game.EndDate,
			}
		}
		e.Image.URL = game.Picture
		payload["embeds"] = append(payload["embeds"], e)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return jsonPayload, nil
}
