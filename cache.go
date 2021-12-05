package cache

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"sync"
	"time"
)

type cache struct {
	dir  string
	lock sync.RWMutex
}
type Data struct {
	Expiration time.Duration `json:"expiration"`
	Unix       int64         `json:"time"`
	Data       []byte        `json:"data"`
}

func (c *cache) Get(key []byte) (data []byte, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var buf []byte
	file := c.getFile(key)
	buf, err = file.Read()
	if err == nil {
		var item Data
		item, err = c.decode(buf)
		if err == nil {
			if item.Expiration == 0 || time.Unix(item.Unix, 0).Add(item.Expiration).After(time.Now()) {
				data = item.Data
			} else {
				err = errors.New("item not in cache")
				defer file.Remove()
			}
		}
	}
	return
}

func (c *cache) Remember(key []byte, expiration time.Duration, fu func() (value []byte, err error)) (data []byte, err error) {
	data, err = c.Get(key)
	if err != nil {
		data, err = fu()
		if err == nil {
			err = c.Put(key, data, expiration)
		}
	}
	return
}

func (c *cache) Put(key []byte, value []byte, expiration time.Duration) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var data []byte
	data, err = c.encode(Data{
		Data:       value,
		Unix:       time.Now().Unix(),
		Expiration: expiration,
	})
	if err == nil {
		err = c.getFile(key).Write(data)
	}
	return
}

func (c *cache) Delete(key []byte) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.getFile(key).Remove()
}

func (c *cache) getFile(key []byte) (file *File) {
	name := md5.Sum(key)
	file = NewFile(filepath.Join(c.dir, "cache", fmt.Sprintf("%x", name[0:1]), fmt.Sprintf("%x", name[1:2]), fmt.Sprintf("%x", name[2:3]), fmt.Sprintf("%x", name[3:])+".db"))
	return
}
func (c *cache) encode(item Data) (data []byte, err error) {
	var t, t2 [8]byte
	i := big.NewInt(int64(item.Expiration)).Bytes()
	i2 := big.NewInt(item.Unix).Bytes()
	copy(t[8-len(i):], i)
	copy(t2[8-len(i2):], i2)
	data = append(data, t[:]...)
	data = append(data, t2[:]...)
	data = append(data, item.Data...)
	return
}

func (c *cache) decode(d []byte) (data Data, err error) {
	data.Expiration = time.Duration(new(big.Int).SetBytes(d[0:8]).Int64())
	data.Unix = new(big.Int).SetBytes(d[8:16]).Int64()
	data.Data = d[16:]
	return
}
func New(dir string) Cache {
	return &cache{dir: dir}
}
