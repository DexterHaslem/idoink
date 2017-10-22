package lastfm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"idoink"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const LastfmCmd = "lastfm"

// 0 = query string, 1 = api key
const lastfmAPI = "http://ws.audioscrobbler.com/2.0/?method=%s&api_key=%s&format=json"

type lastfmApiCreds struct {
	APIKey string `json:"apiKey"`
	Secret string `json:"secret"`
}

type lastfmArtist struct {
	Name      string `json:"name"`
	Playcount string `json:"playcount"` // string in api l o l
}

// TODO: inline anon types for these

type lastfmTopArtists struct {
	Artists []*lastfmArtist `json:"artist"`
}
type lastfmTopArtistsResponse struct {
	TopArtists *lastfmTopArtists `json:"topartists"`
}

var lastfmCreds *lastfmApiCreds

func init() {
	lastfmCreds = &lastfmApiCreds{}
	fb, err := ioutil.ReadFile("lastfm.json")
	if err == nil {
		err = json.Unmarshal(fb, lastfmCreds)
		if err != nil {
			log.Println("failed to load lastfm api key")
		}
	}
}

func lastfmURL(method string) string {
	url := fmt.Sprintf(lastfmAPI, method, lastfmCreds.APIKey)
	return url
}

type lastfmRegisteredInfo struct {
	Text     int64  `json:"#text"`
	UnixTime string `json:"unixtime"`
}

type lastfmUser struct {
	Name       string                `json:"name"`
	RealName   string                `json:"realname"`
	Age        string                `json:"age"`
	Scrobs     string                `json:"playcount"`
	Registered *lastfmRegisteredInfo `json:"registered"`
}

type lastfmUserInfo struct {
	User *lastfmUser `json:"user"`
}

func getUserStats(user string) *lastfmUser {
	method := fmt.Sprintf("user.getinfo&user=%s", user)
	url := lastfmURL(method)

	r, err := http.Get(url)
	if err != nil {
		return nil
	}

	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	ui := &lastfmUserInfo{}
	err = json.Unmarshal(rb, ui)
	if err != nil {
		return nil
	}

	return ui.User
}

func getTopArtists(user string, count int) []*lastfmArtist {
	method := fmt.Sprintf("user.gettopartists&user=%s&limit=10", user)
	url := lastfmURL(method)

	r, err := http.Get(url)
	if err != nil {
		return nil
	}

	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	ta := &lastfmTopArtistsResponse{}
	err = json.Unmarshal(rb, ta)
	if err != nil {
		return nil
	}

	if ta.TopArtists == nil || ta.TopArtists.Artists == nil {
		return nil
	}

	c := count
	got := len(ta.TopArtists.Artists)
	if got == 0 {
		return nil
	}
	if c >= got {
		c = got
	}
	return ta.TopArtists.Artists[0 : c-1]
}

// freaking weirdos

type lastfmTextResponse struct {
	Text string `json:"#text"`
}

type lastfmDate struct {
	// yikes. these guys send EVERYTHING as string
	// well, not everything. funny enough there are a
	// few #text that are NOT STRINGS (user.getinfo).
	// this is the worst api ive seen in a while
	//Unix int64  `json:"uts"`
	Unix string `json:"uts"`
	Text string `json:"#text"`
}

type lastfmPlayedTrack struct {
	// once again not same as artist response
	//Artist *lastfmArtist     `json:"artist"`
	Artist *lastfmTextResponse `json:"artist"`
	Name   string              `json:"name"`
	Album  *lastfmTextResponse `json:"album"`
	Attrs  map[string]string   `json:"@attr"`
	URL    string              `json:"url"`
	Date   *lastfmDate         `json:"date"`
}

type lastfmRecentTracks struct {
	// yes, not plural in api response.. further
	// didnt see value in reusing artist response
	Track []*lastfmPlayedTrack `json:"track"`
}

// argh these bastards
type lastfmRecentTracksResponse struct {
	RecentTracks *lastfmRecentTracks `json:"recenttracks"`
}

func getRecentTracks(user string) []*lastfmPlayedTrack {
	// dont include &nowplaying=true
	// seems to be implicitly on in new api
	method := fmt.Sprintf("user.getrecenttracks&user=%s", user)
	url := lastfmURL(method)

	r, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	rtr := &lastfmRecentTracksResponse{}
	err = json.Unmarshal(rb, rtr)
	if err != nil || rtr.RecentTracks == nil {
		fmt.Println(err)
		return nil
	}
	return rtr.RecentTracks.Track
}

func recentTracksStr(rt []*lastfmPlayedTrack) string {
	// grab first 5, dont care about rest
	b := bytes.Buffer{}
	for i := 0; i < len(rt) && i < 5; i++ {
		nowPlaying := false
		if rt[i].Attrs != nil {
			tnp, ok := rt[i].Attrs["nowplaying"]
			if ok {
				nowPlaying = tnp == "true"
			}
		}
		if nowPlaying {
			b.WriteString(fmt.Sprintf("NOW PLAYING '%s' by %s, ", rt[i].Name, rt[i].Artist.Text))
		} else {
			b.WriteString(fmt.Sprintf("'%s' by %s, ", rt[i].Name, rt[i].Artist.Text))
		}
	}

	return b.String()
}

func artistsStr(ax []*lastfmArtist) string {
	b := bytes.Buffer{}

	for _, a := range ax {
		b.WriteString(fmt.Sprintf("%s (%s), ", a.Name, a.Playcount))
	}
	return b.String()
}

func LastFM(e *idoink.E) (bool, error) { //(from, to string, e.Rest ...string) {
	if lastfmCreds == nil || lastfmCreds.APIKey == "" {
		return false, nil
	}

	if len(e.Rest) < 2 {
		return false, nil
	}

	subcmd := e.Rest[0]
	user := e.Rest[1]

	// TODO: encode
	msg := ""
	switch subcmd {
	case "ta":
		count := 5
		if len(e.Rest) >= 3 {
			tryCount, err := strconv.Atoi(e.Rest[2])
			if err == nil {
				count = tryCount
			}
		}

		ta := getTopArtists(user, count)
		if ta == nil {
			msg = fmt.Sprintf("%s: lastfm - didnt get any top artists for %s", e.From, user)
			break
		}

		artistsStr := artistsStr(ta)

		msg = fmt.Sprintf("%s: lastfm - top artists for %s: %s",
			e.From, user, artistsStr)
		break
	case "u":
		ui := getUserStats(user)

		if ui == nil {
			msg = fmt.Sprintf("%s: lastfm - failed to get user %s", e.From, user)
		} else {
			// convert the timestamp string. hilarious the #text one is not a string, so use it
			rt := time.Unix(ui.Registered.Text, 0)
			msg = fmt.Sprintf("%s: lastfm - user: '%s' registered on %s, name: '%s' age: '%s' total scrobs: '%s'",
				e.From, ui.Name, rt, ui.RealName, ui.Age, ui.Scrobs)
		}

		break
	case "r":
		r := getRecentTracks(user)

		if r == nil {
			msg = fmt.Sprintf("%s: lastfm - failed to get recent tracks for %s", e.From, user)
			break
		}
		recentTrackStr := recentTracksStr(r)

		msg = fmt.Sprintf("%s: lastfm - recent tracks for %s: %s", e.From, user, recentTrackStr)

	default:
		msg = fmt.Sprintf("%s: unknown lastfm command '%s'", e.From, subcmd)
	}

	e.I.Message(e.To, msg)

	return false, nil
}
