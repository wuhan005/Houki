// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"

	"github.com/pkg/errors"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/sqlite"

	"github.com/wuhan005/Houki/assets/migrations"
	"github.com/wuhan005/Houki/internal/dbutil"
)

var Modules ModulesStore

func New() (sqlbuilder.Database, error) {
	db, err := sqlite.Open(sqlite.ConnectionURL{
		Database: "houki.db",
	})
	if err != nil {
		return nil, errors.Wrap(err, "open database")
	}

	_, err = dbutil.Migrate(db.Driver().(*sql.DB), migrations.Migrations)
	if err != nil {
		return nil, errors.Wrap(err, "migrate")
	}

	Modules = NewModulesStore(db)

	return db, nil
}
