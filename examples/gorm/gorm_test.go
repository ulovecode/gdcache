package gorm

import (
	"testing"
)

func TestNewGormCache(t *testing.T) {

	handler := NewGormCacheHandler()

	user := User{
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
	err = handler.GetEntries(&users, "SELECT * FROM user WHERE name = '33'")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users {
		t.Logf("%v", user)
	}

	err = handler.GetEntries(&users, "SELECT * FROM user WHERE id in (3)")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users {
		t.Logf("%v", user)
	}
}
