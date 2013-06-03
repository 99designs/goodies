package ratelimiter

import (
	"bytes"
	"encoding/gob"
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcache struct {
	mc        *memcache.Client
	keyprefix string
}

func NewMemcache(servers string, cacheKeyPrefix string) *Memcache {
	m := new(Memcache)
	m.mc = memcache.New(servers)

	m.keyprefix = cacheKeyPrefix

	return m
}

func (m *Memcache) GetBucketFor(key string) (*LeakyBucket, error) {

	item, err := m.mc.Get(m.keyprefix + key)
	if err != nil {
		return nil, err
	}

	var bucketser LeakyBucketSer
	buf := bytes.NewBuffer(item.Value)
	d := gob.NewDecoder(buf)

	err = d.Decode(&bucketser)
	if err != nil {
		return nil, err
	}

	bucket := bucketser.DeSerialise()

	return bucket, nil
}

func (m *Memcache) SetBucketFor(key string, bucket LeakyBucket) error {

	// LeakyBucket has a closure
	// So we have LeakyBucketSer struct to serialise
	bucketser := bucket.Serialise()

	buf := &bytes.Buffer{}
	e := gob.NewEncoder(buf)
	err := e.Encode(bucketser)

	if err != nil {
		return err
	}

	return m.mc.Set(&memcache.Item{
		Key:        m.keyprefix + key,
		Value:      buf.Bytes(),
		Expiration: int32(bucket.DrainedAt().Unix()),
	})
}
