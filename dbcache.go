package grm

import (
	"database/sql"
	"fmt"
)

var defaultDBCache = NewDBCache()

func RegisterWithDB(alias, driverName, dataSourceName string) (*sql.DB, error) {
	return defaultDBCache.RegisterWithDB(alias, driverName, dataSourceName)
}

func GetWithDB(alias string) (*sql.DB, error) {
	return defaultDBCache.GetWithDB(alias)
}

type DBCache struct {
	mp map[string]*sql.DB
}

func NewDBCache() *DBCache {
	return &DBCache{
		mp: map[string]*sql.DB{},
	}
}

func (d *DBCache) RegisterWithDB(alias, driverName, dataSourceName string) (*sql.DB, error) {
	_, ok := d.mp[alias]
	if ok {
		return nil, fmt.Errorf("%s already exist", alias)
	}
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

func (d *DBCache) GetWithDB(alias string) (*sql.DB, error) {
	db := d.mp[alias]
	if db == nil {
		return nil, fmt.Errorf("%s not exist", alias)
	}
	return db, nil
}
