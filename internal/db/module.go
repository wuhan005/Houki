// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

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

func NewModulesStore(db *gorm.DB) ModulesStore {
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
	*gorm.DB
}

type GetModuleOptions struct {
	EnabledOnly bool
}

func (db *modules) List(ctx context.Context, opts GetModuleOptions) ([]*Module, error) {
	var modules []*Module
	q := db.WithContext(ctx).Model(&Module{})
	if opts.EnabledOnly {
		q = q.Where("enabled = TRUE")
	}
	return modules, q.Find(&modules).Error
}

var ErrModuleNotFound = errors.New("module does not found")

func (db *modules) Get(ctx context.Context, id string) (*Module, error) {
	var module Module
	err := db.WithContext(ctx).Model(&Module{}).Where("id = ?", id).First(&module).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	err := db.WithContext(ctx).Model(&Module{}).Create(&Module{
		ID:   opts.ID,
		Body: opts.Body,
	}).Error
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

	return db.WithContext(ctx).Model(&Module{}).Where("id = ?", id).Set("body = ?", opts.Body).Error
}

func (db *modules) SetStatus(ctx context.Context, id string, enabled bool) error {
	_, err := db.Get(ctx, id)
	if err != nil {
		return err
	}

	return db.WithContext(ctx).Model(&Module{}).Where("id = ?", id).Set("enabled = ?", enabled).Error
}

func (db *modules) Delete(ctx context.Context, id string) error {
	_, err := db.Get(ctx, id)
	if err != nil {
		return err
	}

	return db.WithContext(ctx).Model(&Module{}).Delete("id = ?", id).Error
}
