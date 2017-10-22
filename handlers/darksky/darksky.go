package darksky

import (
	"encoding/json"
	"fmt"
	"idoink"
	"io/ioutil"
	"net/http"
)

const Cmd = "darksky"

type apiKeys struct {
	DarkSkyKey    string `json:"darkSkyKey"`
	ZipCodeApiKey string `json:"zipCodeApiKey"`
}

// api key, zip code
const zipcodeURL = "https://www.zipcodeapi.com/rest/%s/info.json/%s/degrees"

// apikey , lat, long (in degrees)
const darkskyURL = "https://api.darksky.net/forecast/%s/%f,%f"

var keys *apiKeys

func init() {
	keys = &apiKeys{}
	fb, err := ioutil.ReadFile("../apikeys/darksky.json")
	if err == nil {
		json.Unmarshal(fb, keys)
	}
}

func zipToLatLong(z string) (float32, float32, string, error) {
	url := fmt.Sprintf(zipcodeURL, keys.ZipCodeApiKey, z)
	r, err := http.Get(url)
	if err != nil {
		return 0, 0, "", err
	}

	type zr struct {
		Lat  float32 `json:"lat"`
		Lng  float32 `json:"lng"`
		City string  `json:"city"`
		// dont care about rest
	}

	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, 0, "", err
	}
	p := &zr{}
	err = json.Unmarshal(bb, p)
	if err != nil {
		return 0, 0, "", err
	}

	return p.Lat, p.Lng, p.City, nil
}

func DarkSky(e *idoink.E) (bool, error) {
	// we need zip in first cmd
	if len(e.Rest) < 1 {
		return false, nil
	}

	zip := e.Rest[0]
	// convert  zip to lat/long using zipcodeapi first

	lat, long, loc, err := zipToLatLong(zip)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf(darkskyURL, keys.DarkSkyKey, lat, long)
	r, err := http.Get(url)
	if err != nil {
		return false, err
	}

	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}

	// non forecast
	type darkskyResp struct {
		Currently struct {
			TimeUnix            uint64  `json:"time"`
			Summary             string  `json:"summary"`
			Temperature         float32 `json:"temperature"`
			ApparentTemperature float32 `json:"apparentTemperature"`
		} `json:"currently"`
	}

	p := &darkskyResp{}
	err = json.Unmarshal(bb, p)
	if err != nil {
		return false, err
	}

	msg := fmt.Sprintf("%s: it is currently %.1f degrees at %s", e.From, p.Currently.Temperature, loc)
	e.I.Message(e.To, msg)
	return false, nil
}
