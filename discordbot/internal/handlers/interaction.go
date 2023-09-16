package handlers

import (
	"encoding/json"
	"fn/internal/discord"
	"fn/internal/interactions"
	"log"
	"net/http"
)

func Interaction(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)

	msg := make(map[string]interface{})

	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if v, ok := msg["type"].(float64); ok && v == 1 {
		discord.Verify(w, r, body)
		return
	}

	interaction := discord.Interaction{}
	if err := json.Unmarshal(body, &interaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Input\n", string(body))

	var response discord.Response

	switch interaction.Data.Name {
	case "ping":
		response = interactions.Ping(interaction)
	case "courseinfo":
		response = interactions.CourseInfo(interaction)
	default:
		response = discord.Response{
			Type: 4,
			Data: discord.ResponseData{
				Content: "Unknown command",
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Print(err)
	}
}
