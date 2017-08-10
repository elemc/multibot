// -*- Go -*-

package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	log "github.com/sirupsen/logrus"
)

var (
	db *pg.DB
)

// InitDatabase function for initialize pgsql database
func InitDatabase() (err error) {
	var pgo *pg.Options

	if pgo, err = pg.ParseURL(options.PgSQLDSN); err != nil {
		return
	}
	log.Debugf("Try to connect to postgrsql server...")
	db = pg.Connect(pgo)
	err = createTables()
	return
}

func createTables() (err error) {
	tables := []interface{}{}

	for _, t := range tables {
		if err = db.CreateTable(t, &orm.CreateTableOptions{IfNotExists: true}); err != nil {
			return
		}
	}
	return
}
