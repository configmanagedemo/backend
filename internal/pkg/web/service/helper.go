package service

import (
	"net/http"

	"main/internal/pkg/e"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// ResponseOK 正常返回
func ResponseOK(c *gin.Context, data *gin.H) {
	c.JSON(http.StatusOK, gin.H{
		"errcode": e.Success,
		"data":    data,
		"errmsg":  e.GetMsg(e.Success),
	})
	c.Abort()
}

// ResponseError 异常返回
func ResponseError(c *gin.Context, errcode int, errmsg string) {
	if errmsg != "" {
		c.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"data":    "",
			"errmsg":  errmsg,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"errcode": errcode,
			"data":    "",
			"errmsg":  e.GetMsg(errcode),
		})
	}
	c.Abort()
}

// GetUID 从ctx获取uid
func GetUID(c *gin.Context) (uid string, exists bool) {
	if uid, exists := c.Get("uid"); exists {
		return uid.(string), true
	}
	return "", false
}

// GetUUID 获取uuid
func GetUUID() (id string) {
	id = uuid.NewV1().String()
	return
}
