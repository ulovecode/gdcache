package gdcache

import "encoding/json"

// Serializer Serialized abstraction
type Serializer interface {
	Serialize(value interface{}) ([]byte, error)
	Deserialize(data []byte, ptr interface{}) error
}

// JsonSerializer json serializer
type JsonSerializer struct {
}

// Serialize Serialize
func (j JsonSerializer) Serialize(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Deserialize Deserialize
func (j JsonSerializer) Deserialize(data []byte, ptr interface{}) error {
	return json.Unmarshal(data, ptr)
}
