package toodledo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
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

func (c *ToodleClient) Tasks() (*TaskResponse, error) {
	v := url.Values{}
	v.Set("access_token", c.AccessToken)
	v.Set("f", "json")

	resp, err := http.Get(BASE_URL + "/tasks/get.php" + "?" + v.Encode())

	if err != nil {
		return nil, err
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(bytes.NewReader(bts))
	d.UseNumber()
	var results []interface{}
	err = d.Decode(&results)
	log.Print(string(bts))
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
		log.Print(results[0])
		log.Print(reflect.TypeOf(results[0]))
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
