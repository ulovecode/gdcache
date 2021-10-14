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
	mockEntries := make([]MockEntry, 0)
	err = handler.GetEntries(&mockEntries, "SELECT * FROM public_relation where relateId = 1")
	if err != nil {
		t.FailNow()
	}
	for _, m := range mockEntries {
		t.Logf("%v", m)
	}
	err = handler.GetEntries(&mockEntries, "SELECT * FROM public_relation where relateId = 1")
	if err != nil {
		t.FailNow()
	}
	for _, m := range mockEntries {
		t.Logf("%v", m)
	}
	var count int64
	count, err = handler.GetEntriesAndCount(&users1, "SELECT * FROM user WHERE id in (1,2) LIMIT 0,1")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users1 {
		t.Logf("%v", user)
	}
	t.Log(count)

	count, err = handler.GetEntriesAndCount(&users1, "SELECT * FROM user WHERE id in (1,2)")
	if err != nil {
		t.FailNow()
	}
	for _, user := range users1 {
		t.Logf("%v", user)
	}
	t.Log(count)

	users3 := make([]User, 0)
	ids := make([]uint64, 0)
	count, err = handler.GetEntriesAndCount(&users3, "SELECT * FROM user WHERE id in ?", ids)
	if err != nil {
		t.FailNow()
	}
	for _, user := range users1 {
		t.Logf("%v", user)
	}
	t.Log(count)

	count, err = handler.GetEntriesAndCount(&users1, "SELECT * FROM user WHERE id =  ?", nil)
	if err != nil {
		t.FailNow()
	}
	for _, user := range users1 {
		t.Logf("%v", user)
	}
	t.Log(count)

	condition := []User{
		{
			Id: 1,
		},
		{
			Id: 2,
		},
		{
			Id: 3,
		},
	}

	err = handler.GetEntriesByIds(&users1, condition)
	if err != nil {
		t.FailNow()
	}
	for _, user := range users1 {
		t.Logf("%v", user)
	}
	t.Log(count)
}
