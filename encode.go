package gdcache

import (
	"gdcache/schemas"
)

type IEncoder interface {
	encodeIds(ids []schemas.PK) ([]byte, error)
	decodeIds(data []byte) ([]schemas.PK, error)
}

type Encoder struct {
	serializer JsonSerializer
}

func (e Encoder) encodeIds(ids []schemas.PK) ([]byte, error) {
	buf, err := e.serializer.serialize(ids)
	return buf, err
}

func (e Encoder) decodeIds(data []byte) ([]schemas.PK, error) {
	pks := make([]schemas.PK, 0)
	err := e.serializer.deserialize(data, pks)
	return pks, err
}
