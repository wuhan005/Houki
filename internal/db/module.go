// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	dbv3 "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/wuhan005/Houki/internal/dbutil"
	"github.com/wuhan005/Houki/internal/module"
)

var _ ModulesStore = (*modules)(nil)

type ModulesStore interface {
	List(ctx context.Context, opts GetModuleOptions) ([]*Module, error)
	Get(ctx context.Context, id string) (*Module, error)
	Create(ctx context.Context, opts CreateModuleOptions) error
	Update(ctx context.Context, id string, opts UpdateModuleOptions) error
	SetStatus(ctx context.Context, id string, enabled bool) error
	Delete(ctx context.Context, id string) error
}

func NewModulesStore(db sqlbuilder.Database) ModulesStore {
	return &modules{db}
}

// Module represents a single module.
type Module struct {
	ID        string       `db:"id"`
	Body      *module.Body `db:"body"`
	Enabled   bool         `db:"enabled"`
	CreatedAt time.Time    `db:"created_at"`
}

func (m *Module) IsEnabled() bool {
	return m.Enabled
}

type modules struct {
	sqlbuilder.Database
}

type GetModuleOptions struct {
	EnabledOnly bool
}

func (db *modules) List(ctx context.Context, opts GetModuleOptions) ([]*Module, error) {
	var modules []*Module
	q := db.WithContext(ctx).SelectFrom("modules")
	if opts.EnabledOnly {
		q.Where("enabled = TRUE")
	}
	return modules, q.All(&modules)
}

var ErrModuleNotFound = errors.New("module does not found")

func (db *modules) Get(ctx context.Context, id string) (*Module, error) {
	var module Module
	err := db.WithContext(ctx).SelectFrom("modules").Where("id = ?", id).One(&module)
	if err != nil {
		if err == dbv3.ErrNoMoreRows {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}
	return &module, nil
}

type CreateModuleOptions struct {
	ID   string
	Body *module.Body
}

var ErrModuleExists = errors.New("module has already been created")

func (db *modules) Create(ctx context.Context, opts CreateModuleOptions) error {
	_, err := db.WithContext(ctx).InsertInto("modules").
		Columns("id", "body").
		Values(opts.ID, opts.Body).
		Exec()
	if dbutil.IsUniqueViolation(err, "idx_module_id") {
		return ErrModuleExists
	}
	return err
}

type UpdateModuleOptions struct {
	Body *module.Body
}

func (db *modules) Update(ctx context.Context, id string, opts UpdateModuleOptions) error {
	_, err := db.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = db.WithContext(ctx).
		Update("modules").
		Set("body", opts.Body).
		Where("id = ?", id).Exec()
	return err
}

func (db *modules) SetStatus(ctx context.Context, id string, enabled bool) error {
	_, err := db.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = db.WithContext(ctx).
		Update("modules").
		Set("enabled", enabled).
		Where("id = ?", id).Exec()
	return err
}

func (db *modules) Delete(ctx context.Context, id string) error {
	_, err := db.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = db.WithContext(ctx).DeleteFrom("modules").Where("id = ?", id).Exec()
	return err
}
