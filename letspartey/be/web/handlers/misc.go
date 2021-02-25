package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iamlinix.com/partay/web"
)

type WhoAmIResponse struct {
	web.BaseResponse
	Ip string `json:"ip"`
}

func HdlrWhoAmI(c *gin.Context) {
	c.JSON(http.StatusOK, &WhoAmIResponse{
		web.BaseResponse{
			Code:    web.ECOK,
			Message: web.EMOK,
		},
		c.ClientIP(),
	})
}
