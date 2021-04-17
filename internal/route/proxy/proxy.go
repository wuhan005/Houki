package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/wuhan005/gadget"

	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/Houki/internal/proxy"
)

func GetStatus(c *gin.Context) (int, interface{}) {
	return gadget.MakeSuccessJSON(gin.H{
		"enable": proxy.IsEnable(),
	})
}

func Start(c *gin.Context) (int, interface{}) {
	var input struct {
		Address string `json:"address" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		return gadget.MakeErrJSON(40000, "Unexpected proxy address %q: %v", input.Address, err)
	}

	proxy.Start(input.Address)

	if _, err := module.Reload(); err != nil {
		return gadget.MakeErrJSON(50000, "Failed to reload modules: %v", err)
	}
	return gadget.MakeSuccessJSON("success")
}

func Stop(c *gin.Context) (int, interface{}) {
	err := proxy.Stop()
	if err != nil {
		return gadget.MakeErrJSON(40000, "Failed to stop proxy server: %v", err)
	}
	return gadget.MakeSuccessJSON("success")
}
