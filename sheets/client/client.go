package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zinvapel/timetracker/contract"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"strings"
)

type AuthCodeRequest struct {
	Url string
	AuthCode chan string
	MessageId int
}

type Client struct {
	Srv           *sheets.SpreadsheetsService
	SpreadsheetId string
	Err           error
}

func (c Client) Get(readRange string) (*sheets.ValueRange, error) {
	if c.Err != nil {
		return nil, c.Err
	}
	vr, err := c.Srv.Values.Get(c.SpreadsheetId, readRange).Do()

	if err != nil {
		if _, ok := err.(*googleapi.Error); ok || strings.Contains(err.Error(), "oauth2") {
			client = nil
		}
	}

	return vr, err
}

func (c Client) Update(writeRange string, vr *sheets.ValueRange) (err error) {
	if c.Err != nil {
		return c.Err
	}
	_, err = c.Srv.Values.Update(c.SpreadsheetId, writeRange, vr).ValueInputOption("RAW").Do()

	if err != nil {
		if _, ok := err.(*googleapi.Error); ok || strings.Contains(err.Error(), "oauth2") {
			client = nil
		}
	}

	return
}

var client *Client

func GetSheetClient() *Client {
	if client == nil {
		c := obtainClient()
		client = &c
	}

	return client
}

func obtainClient() Client {
	config, err := google.ConfigFromJSON(
		*contract.GetConfig().SheetCredentials,
		"https://www.googleapis.com/auth/spreadsheets",
	)
	if err != nil {
		log.Printf("[sheets.client] Unable to parse client secret file to config: %v\n", err)
		return Client{Err: err}
	}

	srv, err := sheets.New(getClient(config))
	if err != nil {
		log.Printf("[sheets.client] Unable to retrieve Sheets client: %v\n", err)
		return Client{Err: err}
	}

	return Client{Srv: srv.Spreadsheets, SpreadsheetId: *contract.GetConfig().SheetId}
}

func getClient(config *oauth2.Config) *http.Client {
	tokFileName := *contract.GetConfig().SheetTokenFile
	tok, err := tokenFromFile(tokFileName)

	if err != nil {
		tok = getTokenFromWeb(config)
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(tok)

		if tok != nil && err == nil {
			saveToken(*contract.GetConfig().SheetTokenFile, tok)
		}
	}

	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("[sheets.client] Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode = make(chan string)
	contract.GetEventBus().Push(
		contract.NewEvent(
			"google_auth_request",
			&AuthCodeRequest{
				AuthCode: authCode,
				Url: authURL,
			},
			),
		)

	code := <-authCode
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Printf("[sheets.client] Unable to retrieve token from web: %v\n", err)
		return nil
	}

	log.Println("[sheets.client] Token obtained successful")
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("[sheets.client] Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("[sheets.client] Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
