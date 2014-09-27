package toodledo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const BASE_URL = "https://api.toodledo.com/3"

type ToodleClient struct {
	AccessToken string
}

func (c *ToodleClient) AccountInfo() (*Account, error) {
	v := url.Values{}
	v.Set("access_token", c.AccessToken)
	v.Set("f", "json")

	resp, err := http.Get(BASE_URL + "/account/get.php" + "?" + v.Encode())

	if err != nil {
		return nil, err
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var account Account
	err = json.Unmarshal(bts, &account)
	log.Print(string(bts))
	return &account, err
}
