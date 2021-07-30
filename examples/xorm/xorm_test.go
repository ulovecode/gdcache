package xorm

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewXormCache(t *testing.T) {

	var user = User{
		Id: 1,
	}
	handler := NewXormCacheHandler()
	_, err := handler.GetEntry(&user)
	if err != nil {
		panic(err)
	}
	get, _, err := handler.CacheHandler.Get("[id:1]")
	println(string(get))
	v := User{}
	err = json.Unmarshal(get, &v)
	if err!= nil{
		panic(err)
	}
	fmt.Printf("%v", user)

	users := make([]User, 0)
	err = handler.GetEntries(&users, "SELECT * FROM user WHERE id in (1,2)")
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		println(user.Id)
	}
	get, _, err = handler.CacheHandler.Get("[id:1]")
	println(string(get))
	v = User{}
	err = json.Unmarshal(get, &v)
	if err!= nil{
		panic(err)
	}
}
