// +build appengine

package ratelimiter

import (
	"appengine/memcache"
	"bytes"
	"encoding/gob"
)

type Gaecache struct {
	keyprefix string
}

func NewGaecache(keyprefix string) *Gaecache {
	return &Gaecache{
		keyprefix,
	}
}

func (gc *Gaecache) GetBucketFor(key string) (*LeakyBucket, error) {
	if item, err := memcache.Get(gc.keyprefix + key); err != nil {
		return nil, err
	}
	var bucketser LeakyBucketSer
	buf := bytes.NewBuffer(item.Value)
	d := gob.NewDecoder(buf)
	if err := d.Decode(&bucketser); err != nil {
		return nil, err
	}
	bucket := bucketser.DeSerialise()
	return bucket, nil
}

func (gc *Gaecache) SetBucketFor(key string, bucket LeakyBucket) error {
	bucketser := bucket.Serialise()
	buf := &bytes.Buffer{}
	e := gob.NewEncoder(buf)
	if err := e.Encode(bucketser); err != nil {
		return err
	}
	return memcache.Set(&memcache.Item{
		Key:        gc.keyprefix + key,
		Value:      buf.Bytes(),
		Expiration: int32(bucket.DrainedAt().Unix()),
	})
}
