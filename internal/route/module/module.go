package module

import (
	"github.com/gin-gonic/gin"
	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/gadget"
)

func GetModules(c *gin.Context) (int, interface{}) {
	return gadget.MakeSuccessJSON(module.List())
}
