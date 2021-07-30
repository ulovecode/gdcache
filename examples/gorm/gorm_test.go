package gorm

import (
	"fmt"
	"testing"
)

func TestNewGormCache(t *testing.T) {

	var user = User{
		Id: 1,
	}
	_, err := NewGormCacheHandler().GetEntry(&user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", user)

	users := make([]User, 0)
	err = NewGormCacheHandler().GetEntries(&users, "SELECT * FROM user WHERE name = '33'")
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		println(user.Id)
	}
	err = NewGormCacheHandler().GetEntries(&users, "SELECT * FROM user WHERE id in (8)")
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		println(user.Id)
	}
}
