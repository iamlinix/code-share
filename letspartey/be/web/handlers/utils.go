package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"iamlinix.com/partay/json"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/web"
)

func UnmarshalRequestData(c *gin.Context, obj interface{}) error {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("failed to read request data: %v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenCorruptBody,
			Message: web.EMGenCorruptBody,
		})
		return err
	}
	if err = json.JsonUnmarshal(data, obj); err != nil {
		logger.Errorf("error parsing request data: %v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenIncorrectBody,
			Message: web.EMGenIncorrectBody,
		})
		return err
	}

	return nil
}
