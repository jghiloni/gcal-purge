package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/calendar/v3"
)

func getClient(opts *options) {
	l := func(s string, args ...interface{}) {
		if opts.debug != nil {
			opts.debug.Printf(s, args...)
		}
	}

	l("Reading app credentials from $HOME/credentials.json")
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	b, err := ioutil.ReadFile(filepath.Join(home, "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	l("Parsing app credentials")
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokFile := filepath.Join(home, ".gcal-purge")
	tok, err := tokenFromFile(*opts, tokFile)
	if err != nil {
		tok = getTokenFromWeb(*opts, config)
		saveToken(*opts, tokFile, tok)
	}

	opts.client = config.Client(context.Background(), tok)
}

func getTokenFromWeb(opts options, config *oauth2.Config) *oauth2.Token {
	l := func(s string, args ...interface{}) {
		if opts.debug != nil {
			opts.debug.Printf(s, args...)
		}
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	l("Opening %s in a browser", authURL)
	open.Start(authURL)
	fmt.Printf("Go here if it did not automatically launch in your browser: \n\t%s", authURL)

	fmt.Print("\n\nPaste authorization code here: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	l("Getting access token from authorization code")
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(opts options, file string) (*oauth2.Token, error) {
	l := func(s string, args ...interface{}) {
		if opts.debug != nil {
			opts.debug.Printf(s, args...)
		}
	}

	l("Getting token from %s", file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, err
	}
	return tok, err
}

// Saves a token to a file path.
func saveToken(opts options, path string, token *oauth2.Token) {
	l := func(s string, args ...interface{}) {
		if opts.debug != nil {
			opts.debug.Printf(s, args...)
		}
	}

	l("Saving token to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
