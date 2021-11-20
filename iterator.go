package cache

import (
	"time"
)

type Reader interface {
	Get(key []byte) (data []byte, err error)
}
type Writer interface {
	Put(key []byte, value []byte, expiration time.Duration) (err error)
	Delete(key []byte) (err error)
}
type Cache interface {
	Reader
	Writer
	Remember(key []byte, expiration time.Duration, fu func() (value []byte, err error)) (data []byte, err error)
}
