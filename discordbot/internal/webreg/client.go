package webreg

type Client struct {
	cookie string
	term   *Term
}

func NewClient(term *Term) *Client {
	return &Client{
		term: term,
	}
}

func (c *Client) SetCookie(cookie string) {
	c.cookie = cookie
}
