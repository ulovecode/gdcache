package gdcache

type IDB interface {
	// Exec execute sql and return the result
	Exec(sql string, entries interface{}) (err error)
}
