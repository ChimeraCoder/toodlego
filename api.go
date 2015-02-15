package toodlego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Wessie/appdirs"
)

const BASE_URL = "https://api.toodledo.com/3"

var app = appdirs.New("toodlego", "chimeracoder", "0.1")

type ToodleClient struct {
	AppId        string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    float64 `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	Scope        string  `json:"scope"`
	TokenType    string  `json:"token_type"`
}

type ToodleCompleted int

const (
	CompletedOrNot ToodleCompleted = iota
	Uncompleted
	Completed
)

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
	return &account, err
}

func (c *ToodleClient) Tasks(modifiedBefore *time.Time, modifiedAfter *time.Time, completed ToodleCompleted, id int64, start int, fields ...string) (*TaskResponse, error) {
	v := url.Values{}
	v.Set("access_token", c.AccessToken)
	if modifiedBefore != nil {
		v.Set("before", strconv.FormatInt(modifiedBefore.Unix(), 10))
	}
	if modifiedAfter != nil {
		v.Set("after", strconv.FormatInt(modifiedAfter.Unix(), 10))
	}

	if completed != CompletedOrNot {
		// Toodledo uses -1, 0, and 1 instead of 0, 1, 2
		// so we need to subtract 1
		v.Set("comp", strconv.Itoa(int(completed)-1))
	}

	if id != 0 {
		v.Set("id", strconv.FormatInt(id, 10))
	}

	if start != 0 {
		v.Set("start", strconv.Itoa(start))
	}

	v.Set("f", "json")

	v.Set("fields", strings.Join(fields, ","))

	log.Println(v.Encode())
	resp, err := http.Get(BASE_URL + "/tasks/get.php" + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Received response %s", resp.Status)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(bytes.NewReader(bts))
	d.UseNumber()
	var results []interface{}
	err = d.Decode(&results)
	if err != nil {
		log.Print(results)
		return nil, fmt.Errorf("Error decoding JSON response to /tasks/get.php: %s", err)
	}
	if len(results) < 2 {
		return nil, fmt.Errorf("Received only %d responses from /tasks/get.php", len(results))
	}

	taskResponse := &TaskResponse{}

	meta, ok := results[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("First response from /tasks/get.php is not a TaskResponseMeta")
	}
	taskResponse.Meta.Num, _ = meta["num"].(json.Number).Int64()
	taskResponse.Meta.Total, _ = meta["total"].(json.Number).Int64()

	for i, result := range results[1:] {
		bts, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		var task Task
		err = json.Unmarshal(bts, &task)
		if err != nil {
			log.Print(string(bts))
			return nil, fmt.Errorf("Element at index %d of array from /tasks/get.php cannot be unmarshaled into a Task: %s ", i+1, err)
		}

		taskResponse.Tasks = append(taskResponse.Tasks, task)
	}

	return taskResponse, err
}

func (c *ToodleClient) RefreshCredentials() (*RefreshResponse, error) {
	v := url.Values{}
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", c.RefreshToken)

	req, err := http.NewRequest("POST", "https://api.toodledo.com/3/account/token.php", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.AppId, c.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response was: %s", string(bts))

	result := &RefreshResponse{}

	err = json.Unmarshal(bts, result)
	if err != nil {
		return nil, err
	}
	c.AccessToken = result.AccessToken
	c.RefreshToken = result.RefreshToken
	return result, nil
}

func (c *ToodleClient) SaveConfig() error {
	dir := app.UserConfig()
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	filename := path.Join(dir, "config.json")
	bts, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, bts, 0755)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig() (*ToodleClient, error) {
	dir := app.UserConfig()
	filename := path.Join(dir, "config.json")
	var client ToodleClient
	bts, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bts, &client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}
