package autologin

import (
	"context"
	"fmt"
	"log"
	"time"
	"webreg/webreg"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// TODO: fix incomplete logic when login needed but not duo 2 factor

func GetCookies(ctx context.Context, term *webreg.Term, username string, password string) string {
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
		var nodes []*cdp.Node
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

	result := ""
	for _, cookie := range cookies {
		result += fmt.Sprintf("%s=%s;", cookie.Name, cookie.Value)
	}

	return result
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
