package gdcache

import "encoding/json"

type Serializer interface {
	Serialize(value interface{}) ([]byte, error)
	Deserialize(data []byte, ptr interface{}) error
}

type JsonSerializer struct {
}

func (j JsonSerializer) Serialize(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (j JsonSerializer) Deserialize(data []byte, ptr interface{}) error {
	return json.Unmarshal(data, ptr)
}
