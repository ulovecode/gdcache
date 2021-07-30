package gorm

import (
	"testing"
)

func TestNewGormCache(t *testing.T) {
	users := make([]User, 0)
	err := NewGormCacheHandler().GetEntries(&users, "SELECT * FROM user WHERE id in (1,2)")
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		println(user.Id)
	}
}
