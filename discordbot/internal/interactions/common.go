package interactions

import (
	"fmt"
	"fn/internal/discord"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	cookieServer      string
	cookieServerToken string
)

func init() {
	cookieServer = os.Getenv("COOKIE_SERVER")
	cookieServerToken = os.Getenv("COOKIE_SERVER_TOKEN")
}

func discordError(err error) discord.Response {
	response := discord.Response{}
	response.Type = 4
	response.Data.Content = fmt.Sprintf("An error occurred: %s", err.Error())

	return response
}

func getCookie(force bool) (string, error) {
	url := cookieServer
	if force {
		url += "/force"
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth("user", cookieServerToken)

	log.Println("init cookie request", req.URL)
	log.Println("init cookie request", req.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	cookie := string(body)

	return cookie, nil
}
