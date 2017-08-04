package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/gommon/log"
)

const lastfmCmd = "lastfm"

type lastfmApiCreds struct {
	ApiKey string `json:"apiKey"`
	Secret string `json:"secret"`
}

var creds *lastfmApiCreds

func init() {
	creds = &lastfmApiCreds{}
	fb, err := ioutil.ReadFile("lastfmkey.json")
	if err == nil {
		err = json.Unmarshal(fb, creds)
		if err != nil {
			log.Error("failed to load lastfm api key")
		}
	}
}
func lastfm(from, to string, chunks ...string) {
}
