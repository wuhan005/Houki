// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sse

import (
	"sync"
)

type subscriber struct {
	sync.Mutex

	handler      chan *Line
	closeChannel chan struct{}
	closed       bool
}

func (s *subscriber) send(line *Line) {
	select {
	case <-s.closeChannel:
	case s.handler <- line:
	}
}

func (s *subscriber) close() {
	s.Lock()
	if !s.closed {
		close(s.closeChannel)
		s.closed = true
	}
	s.Unlock()
}
