package contract

import (
	"flag"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	ApiToken *string
	Debug *bool
	CuckooFrequency *time.Duration
	// may be personal @todo
	SheetId *string
	UserId *int
	ChatId *int
	SheetTokenFile *string
	SheetCredentials *[]byte
	SheetSchedulePage * string
}

var config Config

func GetConfig() Config {
	return config
}

func init() {
	config = Config{
		ApiToken: flag.String("token", "", "String is api token"),
		Debug: flag.Bool("debug", false, "Enable debug mode"),
		CuckooFrequency: flag.Duration(
			"cuckoo",
			time.Hour,
			"How often generate 'cuckoo' event, format '1d1h2m3s'",
			),
		SheetId: flag.String("sheet", "", "String in https://docs.google.com/spreadsheets/d/<HERE>/edit"),
		UserId: flag.Int("user", 0, "Int user identifier #326793684"),
		ChatId: flag.Int("chat", 0, "Int chat identifier #326793684"),
		SheetTokenFile: flag.String("stoken", "", "Google sheet token.json"),
		SheetSchedulePage: flag.String("ssp", "", "Google sheet schedule page with last ! character"),
	}

	credFile := flag.String("scred", "", "Google sheet credentials.json")

	flag.Parse()

	if *config.ApiToken == "" {
		flag.Usage()
		Finish(NoApiKey)
	}

	f, err := os.Open(*credFile)
	if err != nil {
		flag.Usage()
		Finish(NoCred)
	}

	defer f.Close()

	credBytes, err := ioutil.ReadAll(f)
	if err != nil {
		flag.Usage()
		Finish(NoCred)
	}

	config.SheetCredentials = &credBytes
}