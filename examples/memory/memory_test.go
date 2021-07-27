package main

import (
	"testing"
)

func TestNewMemoryCache(t *testing.T) {
	entries := make([]MockEntry, 0)
	err := NewMemoryCache().GetEntries(&entries, "SELECT * FROM public_relations WHERE id IN (1,2)")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		println(entry.PropertyId)
	}
	entries2 := make([]*MockEntry, 0)
	err = NewMemoryCache().GetEntries(&entries2, "SELECT * FROM public_relations WHERE id IN (1,2)")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries2 {
		println(entry.PropertyId)
	}
}
