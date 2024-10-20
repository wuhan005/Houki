// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

type UpdateCertificate struct {
	Certificate string `json:"certificate" valid:"required"`
	PrivateKey  string `json:"privateKey" valid:"required"`
}
