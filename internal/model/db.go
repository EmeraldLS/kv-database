package model

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var errNotFound = errors.New("no result found")

type Database struct {
	id    string
	name  string
	store map[string]interface{}
	mu    sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		id:    uuid.New().String(),
		name:  "hello_db",
		store: make(map[string]interface{}),
	}
}

func (db *Database) GetId() string {
	return db.id
}

func (db *Database) GetName() string {
	return db.name
}

func (db *Database) GetContent() map[string]interface{} {
	return db.store
}

func (db *Database) WithName(name string) {
	db.name = name
}

func (db *Database) Insert(key string, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.store[key] = value
}

func (db *Database) Find() ([]interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var result = make([]interface{}, 0)
	for _, v := range db.store {
		result = append(result, v)
	}

	return result, nil
}

func (db *Database) FindOne(key string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.store[key] == nil {
		return nil, errNotFound
	}
	result := db.store[key]
	return result, nil
}

func (db *Database) UpdateOne(key string, val interface{}) error {
	_, err := db.FindOne(key)
	if err != nil {
		return err
	}

	db.store[key] = val
	return nil
}

func (db *Database) Remove(key string) error {
	_, err := db.FindOne(key)
	if err != nil {
		return err
	}

	delete(db.store, key)
	return nil
}
