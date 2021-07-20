package gdcache

import "encoding/json"

type Serializer interface {
	serialize(value interface{}) ([]byte, error)
	deserialize(data []byte, ptr interface{}) error
}

type JsonSerializer struct {
}

func (j JsonSerializer) serialize(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (j JsonSerializer) deserialize(data []byte, ptr interface{}) error {
	return json.Unmarshal(data, ptr)
}
