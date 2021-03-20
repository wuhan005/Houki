package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/wuhan005/Houki/internal/proxy"
	"github.com/wuhan005/gadget"
)

func GetStatus(c *gin.Context) (int, interface{}) {
	return gadget.MakeSuccessJSON(gin.H{
		"enable": proxy.IsEnable(),
	})
}

func Start(c *gin.Context) (int, interface{}) {
	proxy.Start()
	return gadget.MakeSuccessJSON("success")
}

func Stop(c *gin.Context) (int, interface{}) {
	err := proxy.Stop()
	if err != nil {
		return gadget.MakeErrJSON(40000, "Failed to stop proxy server: %v", err)
	}
	return gadget.MakeSuccessJSON("success")
}
