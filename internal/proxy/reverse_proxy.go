// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
)

var Reverse = &ReverseProxy{}

// ReverseProxy is the reverse proxy server which supports start and shutdown.
type ReverseProxy struct {
	http.Server

	enable  bool
	address string
}

// IsEnabled returns the proxy status.
func (p *ReverseProxy) IsEnabled() bool {
	return p.enable
}

// Address returns the proxy address.
func (p *ReverseProxy) Address() string {
	return p.address
}
