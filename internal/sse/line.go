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
