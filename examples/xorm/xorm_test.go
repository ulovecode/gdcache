package xorm

import (
	"testing"
)

func TestNewXormCache(t *testing.T) {
	handler := NewXormCacheHandler()

	var user = User{
		Id: 2,
	}
	has, err := handler.GetEntry(&user)
	if err != nil {
		t.FailNow()
	}
	if has {
		t.Logf("%v", user)
	}

	users1 := make([]User, 0)
	err = handler.GetEntries(&users1, "SELECT * FROM user WHERE id in (1,2)")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users1 {
		t.Logf("%v", user)
	}
	t.Logf("1 GetEntries make([]User, 0)")
	users1 = make([]User, 0)
	err = handler.GetEntries(&users1, "SELECT * FROM user WHERE id in (1,2)")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users1 {
		t.Logf("%v", user)
	}
	t.Logf("2 GetEntries make([]User, 0)")
	users2 := make([]*User, 0)
	err = handler.GetEntries(&users2, "SELECT * FROM user WHERE id in (1,2)")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users2 {
		t.Logf("%v", user)
	}
}
