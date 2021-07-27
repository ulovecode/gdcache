package gdcache

type IDB interface {
	// GetEntries cache the entity content obtained through sql, and return the entity of the array pointer type
	GetEntries(entries interface{}, sql string) (interface{}, error)
	// GetEntry get a pointer to an entity type and return the entity
	GetEntry(entry interface{}, sql string) (interface{}, bool, error)
}
