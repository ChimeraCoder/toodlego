package toodledo

import (
	"log"
	"os"
	"reflect"
	"testing"
)

var client = ToodleClient{AppId: os.Getenv("TOODLE_APP_ID"),
	ClientSecret: os.Getenv("TOODLE_CLIENT_SECRET"),
	AccessToken:  os.Getenv("TOODLE_ACCESS_TOKEN"),
	RefreshToken: os.Getenv("TOODLE_REFRESH_TOKEN")}

func TestAccountInfo(t *testing.T) {
	account, err := client.AccountInfo()
	if err != nil {
		t.Error(err)
		return
	}

	if account == nil {
		t.Errorf("Nil account received without error")
		return
	}

	if reflect.DeepEqual(*account, Account{}) {
		t.Errorf("Empty account received")
		return
	}
}

func TestTasks(t *testing.T) {
	taskResponse, err := client.Tasks(nil, nil, Uncompleted, 0, 0, "duedate", "duetime", "startdate", "starttime", "length", "tag", "parent")
	if err != nil {
		t.Error(err)
		return
	}

	if taskResponse == nil {
		t.Errorf("Nil account received without error")
		return
	}

	if reflect.DeepEqual(*taskResponse, TaskResponse{}) {
		t.Errorf("Empty response received for Tasks")
		return
	}

	if len(taskResponse.Tasks) == 0 {
		t.Errorf("Received empty list of tasks")
		return
	}
	for _, task := range taskResponse.Tasks {
		log.Printf("%+v", task)
	}
}

func TestRefresh(t *testing.T) {
	refreshResponse, err := client.RefreshCredentials()
	if err != nil {
		t.Error(err)
		return
	}
	if refreshResponse == nil {
		t.Errorf("Received nil response from RefreshCredentials")
		return
	}
	if reflect.DeepEqual(*refreshResponse, &RefreshResponse{}) {
		t.Errorf("Received empty RefreshResponse")
		return
	}
    log.Printf("Successfully refreshed credentials: %+v", refreshResponse)
}
