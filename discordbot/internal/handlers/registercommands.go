package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	clientId string
	botToken string
)

func init() {
	clientId = os.Getenv("DISCORD_CLIENT_ID")
	botToken = os.Getenv("DISCORD_BOT_TOKEN")
}

func RegisterCommands(w http.ResponseWriter, r *http.Request) {
	if err := registerCommands(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerCommands(w http.ResponseWriter, r *http.Request) error {
	registerURL := fmt.Sprintf("https://discord.com/api/v10/applications/%s/commands", clientId)

	commands := `[
			{
				"name": "ping",
				"description": "Replies with Pong!"
			},
			{
				"name": "courseinfo",
				"description": "Get information about a course",
				"options": [
					{
						"name": "code",
						"description": "The course code",
						"type": 3,
						"required": true
					}
				]
			}
		]`

	req, err := http.NewRequest(http.MethodPut, registerURL, strings.NewReader(commands))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "DiscordBot (https://scottsemtner.com, 1)")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to register commands, status %s", res.Status)
	}

	return nil
}
