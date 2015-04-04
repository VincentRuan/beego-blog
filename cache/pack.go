package cache

import (
	"errors"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// Unmarshal value from memcache.
func Unmarshal(b interface{}, v interface{}) error {
	if b == nil {
		return nil
	}
	bs, ok := b.([]byte)
	if ok {
		return msgpack.Unmarshal(bs, v)
	}
	return errors.New("b must be byte array!")
}
