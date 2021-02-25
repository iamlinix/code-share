package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"iamlinix.com/partay/db"
	"iamlinix.com/partay/fs"
	"iamlinix.com/partay/json"
	"iamlinix.com/partay/logger"
	"iamlinix.com/partay/web"
)

type ReadFile struct {
	Path string `json:"path"`
}

func HdlrPostReadFile(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("failed to read readfile data: %v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenCorruptBody,
			Message: web.EMGenCorruptBody,
		})
		return
	}

	var readFile ReadFile
	if err = json.JsonUnmarshal(data, &readFile); err != nil {
		logger.Errorf("error parsing readfile data: %v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenIncorrectBody,
			Message: web.EMGenIncorrectBody,
		})
		return
	}

	file, err := fs.Get().Open(readFile.Path)
	if err != nil {
		logger.Errorf("error opening file %s: %v", readFile.Path, err)
		c.JSON(http.StatusNotFound, &web.BaseResponse{
			Code:    web.ECResourceNotFound,
			Message: web.EMResourceNotFound,
		})
		return
	}

	data, err = fs.Get().ReadAll(file)
	if err != nil {
		logger.Errorf("error reading file %s: %v", readFile.Path, err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	}

	fs.Get().CloseFile(file)

	c.Data(http.StatusOK, "blob", data)
}

func HdlrReadFile(c *gin.Context) {
	uri, err := url.QueryUnescape(c.Request.RequestURI)
	if err != nil {
		logger.Errorf("error unescaping url %s: %v", c.Request.RequestURI, err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenIncorrectBody,
			Message: web.EMGenIncorrectBody,
		})
		return
	}

	file, err := fs.Get().Open(uri[len("/api/"):])
	if err != nil {
		logger.Errorf("error opening file %s: %v", uri, err)
		c.JSON(http.StatusNotFound, &web.BaseResponse{
			Code:    web.ECResourceNotFound,
			Message: web.EMResourceNotFound,
		})
		return
	}

	data, err := fs.Get().ReadAll(file)
	if err != nil {
		logger.Errorf("error reading file %s: %v", uri, err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	}

	fs.Get().CloseFile(file)

	c.Data(http.StatusOK, "blob", data)
}

type UploadRsp struct {
	web.BaseResponse
	URL string `json:"url"`
}

func HdlrUploadImage(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		logger.Errorf("failed to upload image: %#v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenIncorrectBody,
			Message: web.EMGenIncorrectBody,
		})
		return
	}

	buffer := make([]byte, header.Size)
	file, err := header.Open()
	defer file.Close()

	if err != nil {
		logger.Errorf("failed to open multipart file: %#v", err)
		c.JSON(http.StatusBadRequest, &web.BaseResponse{
			Code:    web.ECGenIncorrectBody,
			Message: web.EMGenIncorrectBody,
		})
		return
	}

	io.ReadFull(file, buffer)

	localPath := fmt.Sprintf("files/%s", uuid.New().String())
	local, err := fs.Get().Create(localPath)
	if err = fs.Get().Write(local, buffer); err != nil {
		logger.Errorf("failed to open local file for image: %#v", err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
		return
	}
	defer fs.Get().CloseFile(local)

	if err = fs.Get().Write(local, buffer); err != nil {
		logger.Errorf("failed to write local image: %#v", err)
		c.JSON(http.StatusInternalServerError, &web.BaseResponse{
			Code:    web.ECGenUnknown,
			Message: web.EMGenUnknown,
		})
	}

	t := c.Request.FormValue("type")
	switch t {
	case "activity":
		id, _ := strconv.Atoi(c.Request.FormValue("lastId"))
		db.Get().Execute("INSERT INTO activity_images (activity_id, url) VALUES (?, ?)",
			id, localPath)
		break

	case "post":
		id, _ := strconv.Atoi(c.Request.FormValue("lastId"))
		db.Get().Execute("INSERT INTO post_images (post_id, url) VALUES (?, ?)",
			id, localPath)
		break

	case "user":
		username := c.Request.FormValue("username")
		db.Get().Execute("UPDATE users SET avatar = ? WHERE username = ?", localPath, username)
		break
	}

	c.JSON(http.StatusOK, &UploadRsp{
		web.BaseResponse{
			Code:    web.ECOK,
			Message: web.EMOK,
		},
		localPath,
	})
}
