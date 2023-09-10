package main

import (
	"flag"
	"webreg/discord"
)

func main() {
	token := flag.String("token", "", "discord bot token")
	flag.Parse()

	discord.Run(*token)
}
