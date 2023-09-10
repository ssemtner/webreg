package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/go-chi/chi"

	"webreg/autologin"
	"webreg/webreg"
)

func main() {
	termFlag := flag.String("term", "", "term code")
	portFlag := flag.Int("port", 3000, "port to listen on")
	usernameFlag := flag.String("username", "", "UCSD username")
	passwordFlag := flag.String("password", "", "UCSD password")

	flag.Parse()

	term, err := webreg.ParseTerm(*termFlag)
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
	r.Get("/cookie", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request received from %s\n", r.RemoteAddr)
		cookies := autologin.GetCookies(ctx, term, username, password)
		log.Printf("Request completed from %s in %s\n", r.RemoteAddr, time.Since(start))

		w.Write([]byte(fmt.Sprintf(`{"cookie":"%s"}`, cookies)))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request received from %s\n", r.RemoteAddr)
		cookies := autologin.GetCookies(ctx, term, username, password)
		log.Printf("Request completed from %s in %s\n", r.RemoteAddr, time.Since(start))

		w.Write([]byte(cookies))
	})

	// warm up call
	autologin.GetCookies(ctx, term, username, password)

	http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), r)
}
