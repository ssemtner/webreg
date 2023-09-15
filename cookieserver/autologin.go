package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// TODO: fix incomplete logic when login needed but not duo 2 factor

func NewContext() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(context.Background())
}

func GetCookies(ctx context.Context, term *Term, username string, password string) string {
	// navigate to webreg
	log.Println("Navigating to webreg...")
	if err := chromedp.Run(ctx, chromedp.Navigate("https://act.ucsd.edu/webreg2/start")); err != nil {
		log.Fatal(err)
	}

	pushSent := false
	i := 0
	for {

		// wait for login, term select, or duo iframe to load (if not already sent)
		selector := "#ssousername,#startpage-select-term"
		if !pushSent {
			selector += ",#duo_iframe"
		}

		log.Println("Waiting for login, term select, or duo iframe to load...")

		if err := chromedp.Run(ctx, chromedp.WaitVisible(selector, chromedp.ByQuery)); err != nil {
			log.Fatal(err)
		}

		// determine which one loaded
		var nodes []*cdp.Node
		if err := chromedp.Run(ctx, chromedp.Nodes(selector, &nodes)); err != nil {
			log.Fatal(err)
		}

		if len(nodes) != 1 {
			log.Fatal("expected one node")
		}

		if nodes[0].AttributeValue("id") == "ssousername" {
			// login form loaded
			log.Println("Login form loaded, logging in...")

			if err := chromedp.Run(ctx, chromedp.Tasks{
				chromedp.SendKeys("#ssousername", username, chromedp.ByID),
				chromedp.SendKeys("#ssopassword", password, chromedp.ByID),
				chromedp.Click("button[type='submit']", chromedp.ByQuery),
			}); err != nil {
				log.Fatal(err)
			}

			log.Println("Logged in, waiting for term select or duo iframe to load...")
		} else if nodes[0].AttributeValue("id") == "duo_iframe" {
			// duo iframe loaded
			log.Println("Duo iframe loaded")

			iframe := nodes[0]

			if err := chromedp.Run(ctx, chromedp.Tasks{
				logAction("Canceling Duo login..."),
				clickInFrame(iframe, ".btn-cancel"),
				chromedp.Sleep(time.Second),

				logAction("Checking remember me..."),
				clickInFrame(iframe, "#remember_me_label_text"),
				chromedp.Sleep(time.Second),

				logAction("Sending another push..."),
				clickInFrame(iframe, "#auth_methods > fieldset > div.row-label.push-label > button"),
				chromedp.Sleep(time.Second),

				logAction("Duo push sent, please approve it"),
			}); err != nil {
				log.Fatal(err)
			}

			pushSent = true

			log.Println("Waiting for term select to load...")
		} else {
			// term select loaded
			log.Printf("Term select loaded, selecting %s\n", term.Code)

			chromedp.Run(ctx, chromedp.Sleep(time.Second))

			if err := chromedp.Run(ctx, chromedp.Tasks{
				chromedp.SetValue("#startpage-select-term", term.Option, chromedp.ByID),
				chromedp.Click("#startpage-button-go", chromedp.ByID),
				chromedp.Sleep(2 * time.Second),
			}); err != nil {
				log.Fatal(err)
			}

			log.Printf("Term %s selected\n", term.Code)

			break
		}

		i += 1
		if i > 10 {
			log.Fatal("Login page, term select page, or duo iframe did not load")
		}
	}

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
