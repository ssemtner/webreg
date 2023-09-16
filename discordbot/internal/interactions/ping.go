package interactions

import "fn/internal/discord"

func Ping(interaction discord.Interaction) discord.Response {
	response := discord.Response{}
	response.Type = 4
	response.Data.Content = "Pong!"

	return response
}
