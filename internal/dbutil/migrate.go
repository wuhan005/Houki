// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"database/sql"
	"embed"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/pkg/errors"
)

func Migrate(db *sql.DB, migrations embed.FS) (func() error, error) {
	var cfg sqlite3.Config
	dbDriver, err := sqlite3.WithInstance(db, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "new database driver")
	}

	srcDriver, err := httpfs.New(http.FS(migrations), ".")
	if err != nil {
		return nil, errors.Wrap(err, "new source driver")
	}

	m, err := migrate.NewWithInstance("httpfs", srcDriver, "sqlite3", dbDriver)
	if err != nil {
		return nil, errors.Wrap(err, "new migrate instance")
	}

	err = m.Up()
	if err == nil || err == migrate.ErrNoChange {
		return dbDriver.Close, nil
	}
	return nil, errors.Wrap(err, "migrate up")
}
