package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	termFlag := flag.String("term", "", "term code")
	portFlag := flag.Int("port", 3000, "port to listen on")
	usernameFlag := flag.String("username", "", "UCSD username")
	passwordFlag := flag.String("password", "", "UCSD password")
	tokenFlag := flag.String("token", "", "token to use for authentication")

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

	// using the basic http auth middleware as a token
	r.Use(middleware.BasicAuth("webreg", map[string]string{
		"user": *tokenFlag,
	}))

	// real ip from nginx
	r.Use(middleware.RealIP)

	// this doubles as the warm up call
	cachedCookies := ""

	// start a warm up call
	go func() {
		cookies := GetCookies(ctx, term, username, password)

		// update the cached cookies
		cachedCookies = strings.Clone(cookies)
	}()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request received from %s\n", r.RemoteAddr)

		// make a copy so they can't change mid request (this is not good still)
		cookies := strings.Clone(cachedCookies)

		// test if the cookies are still valid by trying to get the schedule list
		req, err := http.NewRequest("GET", "https://act.ucsd.edu/webreg2/svc/wradapter/secure/sched-get-schednames?termcode=FA23", nil)
		if err != nil {
			log.Printf("Request failed from %s in %s\n", r.RemoteAddr, time.Since(start))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set("Cookie", cookies)

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Request failed from %s in %s\n", r.RemoteAddr, time.Since(start))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("cookies are invalid, running auth flow...")
			cookies = GetCookies(ctx, term, username, password)

			// update the cached cookies
			cachedCookies = strings.Clone(cookies)
		}

		log.Printf("Request completed from %s in %s\n", r.RemoteAddr, time.Since(start))

		w.Write([]byte(cookies))
	})

	r.Get("/force", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request received (FORCE) from %s\n", r.RemoteAddr)
		cookies := GetCookies(ctx, term, username, password)
		log.Printf("Request completed (FORCE) from %s in %s\n", r.RemoteAddr, time.Since(start))

		w.Write([]byte(cookies))

		cachedCookies = strings.Clone(cookies)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), r)
}
