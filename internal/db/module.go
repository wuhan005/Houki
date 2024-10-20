// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/wuhan005/Houki/internal/dbutil"
	modulespkg "github.com/wuhan005/Houki/internal/modules"
)

var _ ModulesStore = (*modules)(nil)

type ModulesStore interface {
	List(ctx context.Context, opts ListModuleOptions) ([]*Module, int64, error)
	All(ctx context.Context, opts AllModuleOptions) ([]*Module, error)
	Get(ctx context.Context, id uint) (*Module, error)
	Create(ctx context.Context, opts CreateModuleOptions) (*Module, error)
	Update(ctx context.Context, id uint, opts UpdateModuleOptions) error
	SetStatus(ctx context.Context, id uint, enabled bool) error
	Delete(ctx context.Context, id uint) error
}

func NewModulesStore(db *gorm.DB) ModulesStore {
	return &modules{db}
}

// Module represents a single module.
type Module struct {
	dbutil.Model

	Name    string           `json:"name"`
	Body    *modulespkg.Body `json:"body"`
	Enabled bool             `json:"enabled"`
}

func (m *Module) IsEnabled() bool {
	return m.Enabled
}

type modules struct {
	*gorm.DB
}

type ListModuleOptions struct {
	dbutil.Pagination
	EnabledOnly bool
}

func (db *modules) List(ctx context.Context, opts ListModuleOptions) ([]*Module, int64, error) {
	var modules []*Module
	q := db.WithContext(ctx).Model(&Module{})
	if opts.EnabledOnly {
		q = q.Where("enabled = TRUE")
	}

	var count int64
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrap(err, "count")
	}

	limit, offset := dbutil.LimitOffset(opts.Page, opts.PageSize)
	q = q.Limit(limit).Offset(offset)
	if err := q.Find(&modules).Error; err != nil {
		return nil, 0, errors.Wrap(err, "find")
	}

	return modules, count, nil
}

type AllModuleOptions struct {
	EnabledOnly bool
}

func (db *modules) All(ctx context.Context, opts AllModuleOptions) ([]*Module, error) {
	q := db.WithContext(ctx).Model(&Module{})
	if opts.EnabledOnly {
		q = q.Where("enabled = TRUE")
	}

	var modules []*Module
	if err := q.Find(&modules).Error; err != nil {
		return nil, errors.Wrap(err, "find")
	}
	return modules, nil
}

var ErrModuleNotFound = errors.New("module does not found")

func (db *modules) Get(ctx context.Context, id uint) (*Module, error) {
	var module Module
	if err := db.WithContext(ctx).Model(&Module{}).Where("id = ?", id).First(&module).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModuleNotFound
		}
		return nil, err
	}
	return &module, nil
}

type CreateModuleOptions struct {
	Name string
	Body *modulespkg.Body
}

func (db *modules) Create(ctx context.Context, opts CreateModuleOptions) (*Module, error) {
	module := &Module{
		Name: opts.Name,
		Body: opts.Body,
	}
	if err := db.WithContext(ctx).Model(&Module{}).Create(&module).Error; err != nil {
		return nil, err
	}
	return module, nil
}

type UpdateModuleOptions struct {
	Name string
	Body *modulespkg.Body
}

func (db *modules) Update(ctx context.Context, id uint, opts UpdateModuleOptions) error {
	return db.WithContext(ctx).Model(&Module{}).Where("id = ?", id).
		Update("name", opts.Name).
		Update("body", opts.Body).
		Error
}

func (db *modules) SetStatus(ctx context.Context, id uint, enabled bool) error {
	return db.WithContext(ctx).Model(&Module{}).Where("id = ?", id).Set("enabled", enabled).Error
}

func (db *modules) Delete(ctx context.Context, id uint) error {
	return db.WithContext(ctx).Model(&Module{}).Delete("id = ?", id).Error
}
