package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseHelper struct {
	Status  uint        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BadResponse(c *gin.Context, msg string) {
	res := responseHelper{Status: http.StatusBadRequest, Data: nil, Message: msg}
	c.JSON(http.StatusBadRequest, res)
}

func SuccessResponse(c *gin.Context, msg string, data interface{}) {
	res := responseHelper{Status: http.StatusOK, Data: data, Message: msg}
	c.JSON(http.StatusOK, res)
}
