package db

import "github.com/ulovecode/gdcache/schemas"

type IDB interface {
	// GetEntries cache the entity content obtained through sql, and return the entity of the array pointer type
	GetEntries(entries []schemas.IEntry, sql string) error
	// GetEntry get a pointer to an entity type and return the entity
	GetEntry(entry schemas.IEntry, sql string) (bool, error)
}
