// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package module

import (
	"github.com/gin-gonic/gin"
	"github.com/wuhan005/gadget"

	"github.com/wuhan005/Houki/internal/module"
)

func ListModules(c *gin.Context) (int, interface{}) {
	moduleLists, err := module.Scan()
	if err != nil {
		return gadget.MakeErrJSON(50000, "Failed to list modules: %v", err)
	}
	enabledModules, err := module.GetEnabledModules()
	if err != nil {
		return gadget.MakeErrJSON(50000, "Failed to get enabled modules: %v", err)
	}
	return gadget.MakeSuccessJSON(gin.H{
		"list":            moduleLists,
		"enabled_modules": enabledModules,
	})
}

func EnableModule(c *gin.Context) (int, interface{}) {
	modID := c.Param("id")

	if err := module.Enable(modID); err != nil {
		return gadget.MakeErrJSON(50000, "Failed to enable module: %v", err)
	}
	return gadget.MakeSuccessJSON("Enable module succeed!")
}

func DisableModule(c *gin.Context) (int, interface{}) {
	modID := c.Param("id")

	if err := module.Disable(modID); err != nil {
		return gadget.MakeErrJSON(50000, "Failed to disable module: %v", err)
	}
	return gadget.MakeSuccessJSON("Disable module succeed!")
}
