package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

type tokenInfo struct {
	Token  string
	Expire time.Time
}

// Store is store for token and expires.
type Store interface {
	Save(token string, expire time.Time) error
	Token() string
	Expired() bool
}

type fileStore struct {
	tokenInfo tokenInfo
	path      string
}

// NewFileStore returns Store from path of credential file
func NewFileStore(path string) (Store, error) {
	var ti tokenInfo
	c, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("fail to open credential store: %w", err)
	}
	if !os.IsNotExist(err) {
		dec := gob.NewDecoder(c)
		if err := dec.Decode(&ti); err != nil {
			return nil, fmt.Errorf("fail to decode credential file: %w", err)
		}
	}

	fs := &fileStore{tokenInfo: ti, path: path}

	return fs, nil
}

func (f *fileStore) Save(token string, expire time.Time) error {
	c, err := os.Create(f.path)
	if err != nil {
		return fmt.Errorf("fail to create credential file: %w", err)
	}
	enc := gob.NewEncoder(c)
	err = enc.Encode(tokenInfo{Token: token, Expire: expire})
	if err != nil {
		return fmt.Errorf("fail to encode credential: %w", err)
	}
	return nil
}

func (f *fileStore) Token() string { return f.tokenInfo.Token }
func (f *fileStore) Expired() bool { return f.tokenInfo.Expire.Before(time.Now()) }
