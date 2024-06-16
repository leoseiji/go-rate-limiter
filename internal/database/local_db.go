package database

import (
	"errors"
	"sync"
	"time"
)

type user struct {
	Counter int64
	Ttl     time.Time
}

type LocalDB struct {
	mu      *sync.RWMutex
	storage map[string]user
}

func NewLocalDB() *LocalDB {
	db := &LocalDB{
		mu:      &sync.RWMutex{},
		storage: make(map[string]user),
	}
	go db.cleanUp()
	return db
}

func (db *LocalDB) Get(key string) (int64, bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, exists := db.storage[key]
	if !exists {
		return 0, false, nil
	}
	return user.Counter, exists, nil
}

func (db *LocalDB) Set(key string, value int64, expiration time.Duration) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user := user{
		Counter: value,
		Ttl:     time.Now().Add(expiration),
	}
	db.storage[key] = user

	return nil
}

func (db *LocalDB) Incr(key string) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, exists := db.storage[key]
	if !exists {
		return 0, errors.New("user not found")
	}
	user.Counter++
	db.storage[key] = user

	return user.Counter, nil
}

func (db *LocalDB) cleanUp() {
	for {
		for {
			time.Sleep(time.Second * 1)

			db.mu.RLock()
			for key, user := range db.storage {
				if time.Now().After(user.Ttl) {
					delete(db.storage, key)
				}
			}
			db.mu.RUnlock()
		}
	}
}
