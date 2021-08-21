// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbutil

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func isConstraintViolation(err error, code pq.ErrorCode, constraint string) bool {
	pgErr, ok := errors.Cause(err).(*pq.Error)
	if !ok {
		return false
	}
	return pgErr.Code == code && pgErr.Constraint == constraint
}

// IsInvalidUUID checks the UUID format error.
func IsInvalidUUID(err error) bool {
	return isConstraintViolation(err, "22P02", "")
}

// IsUniqueViolation checks the unique violation error.
func IsUniqueViolation(err error, constraint string) bool {
	return isConstraintViolation(err, "23505", constraint)
}

// IsForeignKeyViolation checks the foreign key violation error.
func IsForeignKeyViolation(err error, constraint string) bool {
	return isConstraintViolation(err, "23503", constraint)
}

// IsCheckViolation checks the value validation violation error.
func IsCheckViolation(err error, constraint string) bool {
	return isConstraintViolation(err, "23514", constraint)
}
