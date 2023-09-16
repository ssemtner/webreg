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
	TTS     bool             `json:"tts"`
	Content string           `json:"content"`
	Embeds  []*ResponseEmbed `json:"embeds"`
}

// An MessageEmbed stores data for message embeds.
type ResponseEmbed struct {
	URL         string                  `json:"url,omitempty"`
	Type        EmbedType               `json:"type,omitempty"`
	Title       string                  `json:"title,omitempty"`
	Description string                  `json:"description,omitempty"`
	Timestamp   string                  `json:"timestamp,omitempty"`
	Color       int                     `json:"color,omitempty"`
	Footer      *ResponseEmbedFooter    `json:"footer,omitempty"`
	Image       *ResponseEmbedImage     `json:"image,omitempty"`
	Thumbnail   *ResponseEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *ResponseEmbedVideo     `json:"video,omitempty"`
	Provider    *ResponseEmbedProvider  `json:"provider,omitempty"`
	Author      *ResponseEmbedAuthor    `json:"author,omitempty"`
	Fields      []*ResponseEmbedField   `json:"fields,omitempty"`
}

// ResponseEmbedThumbnail is a part of a MessageEmbed struct.
type ResponseEmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// ResponseEmbedVideo is a part of a MessageEmbed struct.
type ResponseEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// ResponseEmbedProvider is a part of a MessageEmbed struct.
type ResponseEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// ResponseEmbedAuthor is a part of a MessageEmbed struct.
type ResponseEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// ResponseEmbedField is a part of a MessageEmbed struct.
type ResponseEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// ResponseEmbedFooter is a part of a MessageEmbed struct.
type ResponseEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// ResponseEmbedImage is a part of a MessageEmbed struct.
type ResponseEmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// https://discord.com/developers/docs/resources/channel#embed-object-embed-types
type EmbedType string

// Block of valid EmbedTypes
const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGifv    EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)
