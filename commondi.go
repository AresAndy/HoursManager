package main

import (
	"database/sql"

	"github.com/sarulabs/di"
)

func createDIbuilder() (*di.Builder, error) {
	builder, err := di.NewBuilder()

	if err != nil {
		return nil, err
	}

	builder.Add(di.Def{
		Name: "db",
		Build: func(ctn di.Container) (interface{}, error) {
			sqliteFilePath := "/home/silvio/hours.sqlite"
			return sql.Open("sqlite3", "file:"+sqliteFilePath)
		},
		Close: func(obj interface{}) error {
			return obj.(*sql.DB).Close()
		},
	})

	return builder, err
}
