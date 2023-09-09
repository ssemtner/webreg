package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func main() {
	termFlag := flag.String("term", "", "term code")
	portFlag := flag.Int("port", 3000, "port to listen on")
	usernameFlag := flag.String("username", "", "UCSD username")
	passwordFlag := flag.String("password", "", "UCSD password")

	flag.Parse()

	term, err := ParseTerm(*termFlag)
	if err != nil {
		log.Fatal(err)
	}

	username := *usernameFlag
	password := *passwordFlag

	if username == "" {
		username = os.Getenv("WEBREG_USERNAME")
	}

	if password == "" {
		password = os.Getenv("WEBREG_PASSWORD")
	}

	if username == "" || password == "" {
		log.Fatal("please provide a username and password via flags or the WEBREG_USERNAME and WEBREG_PASSWORD environment variables")
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		cookies := getCookies(ctx, term, username, password)

		list := []render.Renderer{}
		for _, cookie := range cookies {
			list = append(list, NewCookieResponse(cookie))
		}

		render.RenderList(w, r, list)

		log.Println("Request completed")
	})
	http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), r)
}

func getCookies(ctx context.Context, term *Term, username string, password string) []*network.Cookie {
	// navigate to webreg
	log.Println("Navigating to webreg...")
	if err := chromedp.Run(ctx, chromedp.Navigate("https://act.ucsd.edu/webreg2/start")); err != nil {
		log.Fatal(err)
	}

	// wait for term select or login form to load
	log.Println("Waiting for term select or login form to load...")
	if err := chromedp.Run(ctx, chromedp.WaitVisible("#startpage-button-go,#ssousername", chromedp.ByID)); err != nil {
		log.Fatal(err)
	}

	// determine which one loaded
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes("#startpage-button-go,#ssousername", &nodes)); err != nil {
		log.Fatal(err)
	}

	if len(nodes) != 1 {
		log.Fatal("expected one node")
	}

	// run login procedure if login form loaded
	if nodes[0].AttributeValue("id") == "ssousername" {
		log.Println("Logging in...")

		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.SendKeys("#ssousername", username, chromedp.ByID),
			chromedp.SendKeys("#ssopassword", password, chromedp.ByID),
			chromedp.Click("button[type='submit']", chromedp.ByQuery),
		}); err != nil {
			log.Fatal(err)
		}

		log.Println("Waiting for duo iframe or go button to load...")

		// wait for duo iframe OR go button to load
		if err := chromedp.Run(ctx, chromedp.WaitVisible("#duo_iframe,#startpage-button-go", chromedp.ByQuery)); err != nil {
			log.Fatal(err)
		}

		// determine which one loaded
		if err := chromedp.Run(ctx, chromedp.Nodes("#duo_iframe,#startpage-button-go", &nodes)); err != nil {
			log.Fatal(err)
		}

		if len(nodes) != 1 {
			log.Fatal("expected one node")
		}

		// run duo login procedure if duo iframe loaded
		if nodes[0].AttributeValue("id") == "duo_iframe" {
			iframe := nodes[0]

			if err := chromedp.Run(ctx, chromedp.Tasks{
				logAction("Canceling Duo login..."),
				clickInFrame(iframe, ".btn-cancel"),
				chromedp.Sleep(2 * time.Second),

				logAction("Checking remember me..."),
				clickInFrame(iframe, "#remember_me_label_text"),
				chromedp.Sleep(2 * time.Second),

				logAction("Sending another push..."),
				clickInFrame(iframe, "#auth_methods > fieldset > div.row-label.push-label > button"),
				chromedp.Sleep(2 * time.Second),

				logAction("Waiting for term selection page..."),
				chromedp.WaitVisible("#startpage-button-go", chromedp.ByID),
			}); err != nil {
				log.Fatal(err)
			}

			log.Println("Term selection page loaded")
		} else {
			log.Println("Go button loaded, skipping Duo login...")
		}
	}

	// select term
	log.Println("Selecting term...")

	chromedp.Run(ctx, chromedp.Sleep(time.Second))

	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.SetValue("#startpage-select-term", term.Option, chromedp.ByID),
		chromedp.Click("#startpage-button-go", chromedp.ByID),
		chromedp.Sleep(2 * time.Second),
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("Term %s selected", term.Code)

	// extract cookies
	log.Println("Extracting cookies...")
	var cookies []*network.Cookie
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate("https://act.ucsd.edu/webreg2/svc/wradapter/secure/sched-get-schednames?termcode=FA23"),
		chromedp.Sleep(3 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			c, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			cookies = c
			return nil
		}),
	}); err != nil {
		log.Fatal(err)
	}

	return cookies
}

type Term struct {
	Code   string
	Option string
}

func ParseTerm(code string) (*Term, error) {
	options := map[string]string{
		"FA23": "5320:::FA23",
	}

	option, ok := options[code]
	if !ok {
		return nil, fmt.Errorf("invalid term code")
	}

	return &Term{
		Code:   code,
		Option: option,
	}, nil
}

func logAction(value string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Println(value)

		return nil
	})
}

func clickInFrame(iframe *cdp.Node, selector string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		var nodes []*cdp.Node

		if err := chromedp.Nodes(selector, &nodes).Do(ctx); err != nil {
			return err
		}

		return chromedp.MouseClickNode(nodes[0]).Do(ctx)
	}
}

type CookieResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Elapsed int64 `json:"elapsed"`
}

func NewCookieResponse(cookie *network.Cookie) *CookieResponse {
	return &CookieResponse{
		Name:  cookie.Name,
		Value: cookie.Value,
	}
}

func (rd *CookieResponse) Render(w http.ResponseWriter, r *http.Request) error {
	rd.Elapsed = 10
	return nil
}