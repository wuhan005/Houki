// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sse

import (
	"time"
)

// Line is a single line of the log.
type Line struct {
	Type      string      `json:"type"`
	Message   interface{} `json:"message"`
	Timestamp int64       `json:"time"`
}

// NewLine creates a line.
func NewLine(messageType string, message interface{}) *Line {
	return &Line{
		Type:      messageType,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
}
