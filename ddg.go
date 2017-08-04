package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const ddgCmd = "ddg"

const ddgAPI = "https://api.duckduckgo.com/?q=%s&format=json&t=dmhbot"

type ddgResult struct {
	Result   string `json:"Result"`
	FirstURL string `json:"FirstURL"`
	Text     string `json:"Text"`
}
type ddgResp struct {
	Abstract       string       `json:"Abstract"`
	AbstractText   string       `json:"AbstractText"`
	AbstractSource string       `json:"AbstractSource"`
	AbstractURL    string       `json:"AbstractURL"`
	Image          string       `json:"Image"`
	Heading        string       `json:"Heading"`
	Answer         string       `json:"Answer"`
	AnswerType     string       `json:"AnswerType"`
	Results        []*ddgResult `json:"Results"`
	Type           string       `json:"Type"`
}

func ddg(from, to string, chunks ...string) {
	// squish message to flat string and then query it
	q := strings.Join(chunks, "")
	url := fmt.Sprintf(ddgAPI, q)

	go func() {
		r, err := http.Get(url)
		if err != nil {
			i.PrivMsg(to, "ddg: i messed up requesting to search ddg :-(")
			return
		}

		rb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			i.PrivMsg(to, "ddg: i messed up reading response from ddg :-(")
		}

		pr := &ddgResp{}
		err = json.Unmarshal(rb, pr)
		if err != nil {
			i.PrivMsg(to, "ddg: i messed up reading response from ddg :-(")
		}
		msg := ""
		if len(pr.Results) > 0 {
			u := pr.Results[0].FirstURL
			t := pr.Results[0].Text
			msg = fmt.Sprintf("%s: ddg - %s - %s (%s)", from, t, u, pr.Type)
		} else {
			msg = fmt.Sprintf("%s: ddg - I didnt find any instant answers (the api sucks)", from)
		}
		i.PrivMsg(to, msg)
	}()
}
