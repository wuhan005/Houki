// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"encoding/json"
)

type NewModule struct {
	ID   string
	Body json.RawMessage
}

type UpdateModule struct {
	Body json.RawMessage
}