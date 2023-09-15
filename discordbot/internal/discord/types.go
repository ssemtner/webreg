package discord

// Interaction is the request we get from Discord when a user
// triggers a slash Command i.e. /zoom
type Interaction struct {
	Type   float64           `json:"type"`
	Data   InteractionData   `json:"data"`
	Member InteractionMember `json:"member"`
}

// InteractionData is present for the slash command itself
// i.e. /zoom
type InteractionData struct {
	Name    string                   `json:"name"`
	ID      string                   `json:"id"`
	Type    float64                  `json:"type"`
	Options []InteractionDataOptions `json:"options"`
}

// InteractionDataOptions contains the option passed in
// within the slash command i.e. the parameters
type InteractionDataOptions struct {
	Name  string      `json:"name"`
	Type  float64     `json:"type"`
	Value interface{} `json:"value"`
}

type InteractionMember struct {
	User InteractionMemberUser `json:"user"`
}

// InteractionMemberUser gives a way to uniquely
// identify a user by adding # between the Username and
// the Discriminator
type InteractionMemberUser struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

// Response is the response we send back to Discord
// See also: https://discord.com/developers/docs/interactions/receiving-and-responding
type Response struct {
	Type float64      `json:"type"`
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Content string          `json:"content"`
	Embeds  []ResponseEmbed `json:"embeds,omitempty"`
}

type ResponseEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Type        string `json:"type"`
}
