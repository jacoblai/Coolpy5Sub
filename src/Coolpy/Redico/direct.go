package Redico

import "errors"

var (
	// ErrKeyNotFound is returned when a key doesn't exist.
	ErrKeyNotFound = errors.New(msgKeyNotFound)
	// ErrWrongType when a key is not the right type.
	ErrWrongType = errors.New(msgWrongType)
	// ErrIntValueError can returned by INCRBY
	ErrIntValueError = errors.New(msgInvalidInt)
	// ErrFloatValueError can returned by INCRBYFLOAT
	ErrFloatValueError = errors.New(msgInvalidFloat)
)

// Select sets the DB id for all direct commands.
func (m *Redico) Select(i int) {
	m.Lock()
	defer m.Unlock()
	m.selectedDB = i
}

// Get returns string keys added with SET.
func (m *Redico) Get(k string) (string, error) {
	return m.DB(m.selectedDB).Get(k)
}

// Get returns a string key
func (db *RedicoDB) Get(k string) (string, error) {
	db.master.Lock()
	defer db.master.Unlock()
	if !db.exists(k) {
		return "", ErrKeyNotFound
	}
	return db.stringGet(k), nil
}

// Set sets a string key. Removes expire.
func (m *Redico) Set(k, v string) error {
	return m.DB(m.selectedDB).Set(k, v)
}

// Set sets a string key. Removes expire.
// Unlike redis the key can't be an existing non-string key.
func (db *RedicoDB) Set(k, v string) error {
	db.master.Lock()
	defer db.master.Unlock()

	if db.exists(k) {
		return ErrKeyNotFound
	}
	db.del(k, true) // Remove expire
	db.stringSet(k, v)
	return nil
}

// Del deletes a key and any expiration value. Returns whether there was a key.
func (m *Redico) Del(k string) bool {
	return m.DB(m.selectedDB).Del(k)
}

// Del deletes a key and any expiration value. Returns whether there was a key.
func (db *RedicoDB) Del(k string) bool {
	db.master.Lock()
	defer db.master.Unlock()
	if !db.exists(k) {
		return false
	}
	db.del(k, true)
	return true
}

// Exists tells whether a key exists.
func (m *Redico) Exists(k string) bool {
	return m.DB(m.selectedDB).Exists(k)
}

// Exists tells whether a key exists.
func (db *RedicoDB) Exists(k string) bool {
	db.master.Lock()
	defer db.master.Unlock()
	return db.exists(k)
}

