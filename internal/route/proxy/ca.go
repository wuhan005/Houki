// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wuhan005/gadget"

	"github.com/wuhan005/Houki/internal/ca"
)

func GenerateCA(c *gin.Context) (int, interface{}) {
	crtBytes, _, err := ca.Generate(true)
	if err != nil {
		return gadget.MakeErrJSON(50000, "Failed to generate CA: %v", err)
	}
	return gadget.MakeSuccessJSON(string(crtBytes))
}

func FetchCA(c *gin.Context) (int, interface{}) {
	crtBytes, err := os.ReadFile(".certificate/ca.crt")
	if err != nil {
		return gadget.MakeErrJSON(40400, "File .certificate/ca.crt not found: %v", err)
	}

	return gadget.MakeSuccessJSON(string(crtBytes))
}
