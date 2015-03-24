package cache

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/cache"
	"github.com/bradfitz/gomemcache/memcache"
	"gopkg.in/vmihailenco/msgpack.v2"
	"strings"
)

type PackMemcacheCache struct {
	conn     *memcache.Client
	conninfo []string
}

// create new memcache adapter.
func NewPackMemCache() *PackMemcacheCache {
	return &PackMemcacheCache{}
}

// Unmarshal value from memcache.
func Unmarshal(b interface{}, v interface{}) error {
	bs, ok := b.([]byte)
	if ok {
		return msgpack.Unmarshal(bs, v)
	}
	return errors.New("b must be byte array!")
}

// get value from memcache.
func (rc *PackMemcacheCache) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return string(item.Value)
	}
	return nil
}

// put value to memcache.
func (rc *PackMemcacheCache) Put(key string, val interface{}, timeout int64) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}

	b, err := msgpack.Marshal(val)
	if err != nil {
		return err
	}
	item := memcache.Item{Key: key, Value: b, Expiration: int32(timeout)}
	return rc.conn.Set(&item)
}

// delete value in memcache.
func (rc *PackMemcacheCache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// increase counter.
func (rc *PackMemcacheCache) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// decrease counter.
func (rc *PackMemcacheCache) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// check value exists in memcache.
func (rc *PackMemcacheCache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	_, err := rc.conn.Get(key)
	if err != nil {
		return false
	}
	return true
}

// clear all cached in memcache.
func (rc *PackMemcacheCache) ClearAll() error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.FlushAll()
}

// start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *PackMemcacheCache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *PackMemcacheCache) connectInit() error {
	rc.conn = memcache.New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("packmemcache", NewPackMemCache())
}
