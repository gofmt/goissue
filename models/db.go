package models

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	db     *pg.DB
	tables = []interface{}{
		(*User)(nil),
		(*Issue)(nil),
	}
)

func Connect(ctx context.Context, dsn string) error {
	opt, err := pg.ParseURL(dsn)
	if err != nil {
		return err
	}

	db = pg.Connect(opt).WithContext(ctx)
	for _, table := range tables {
		if err := db.CreateTable(table, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		}); err != nil {
			return err
		}
	}

	return nil
}
