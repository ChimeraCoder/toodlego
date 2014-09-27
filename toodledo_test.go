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
