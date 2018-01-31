package grm

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
)

var defaultDBCache = NewDBCache()

func RegisterWithDB(alias, source string) (*sql.DB, error) {
	return defaultDBCache.RegisterWithDB(alias, source)
}

func Register(source string) (*sql.DB, error) {
	return defaultDBCache.Register(source)
}

func GetWithDB(alias string) (*sql.DB, error) {
	return defaultDBCache.GetWithDB(alias)
}

func Get() (*sql.DB, error) {
	return defaultDBCache.Get()
}

type DBCache struct {
	mp map[string]*sql.DB
}

func NewDBCache() *DBCache {
	return &DBCache{
		mp: map[string]*sql.DB{},
	}
}

func (d *DBCache) RegisterWithDB(alias, source string) (*sql.DB, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("source")
	}

	_, ok := d.mp[alias]
	if ok {
		return nil, fmt.Errorf("%s already exist", alias)
	}

	driverName := ""
	dataSourceName := ""
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}
	if u.Opaque != "" {
		return nil, errors.New("Error RegisterWithDB: " + source)
	}
	driverName = u.Scheme
	dataSourceName = source[len(u.Scheme)+3:]

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	d.mp[alias] = db
	return db, nil
}

func (d *DBCache) Register(source string) (*sql.DB, error) {
	return d.RegisterWithDB("", source)
}

func (d *DBCache) GetWithDB(alias string) (*sql.DB, error) {
	db := d.mp[alias]
	if db == nil {
		return nil, fmt.Errorf("%s not exist", alias)
	}
	return db, nil
}

func (d *DBCache) Get() (*sql.DB, error) {
	return d.GetWithDB("")
}
