package xorm

import (
	"testing"
)

func TestNewXormCache(t *testing.T) {
	handler := NewXormCacheHandler()

	var user = User{
		Id: 1,
	}
	has, err := handler.GetEntry(&user)
	if err != nil {
		t.FailNow()
	}
	if has {
		t.Logf("%v", user)
	}

	users := make([]User, 0)
	err = handler.GetEntries(&users, "SELECT * FROM user WHERE id in (1,2)")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users {
		t.Logf("%v", user)
	}
}
