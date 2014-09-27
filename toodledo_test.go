package toodledo

import (
	"os"
	"reflect"
	"testing"
)

var client = ToodleClient{os.Getenv("TOODLE_ACCESS_TOKEN")}

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
	taskResponse, err := client.Tasks("duedate", "duetime")
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
}
