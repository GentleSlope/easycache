package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/pkg/upload"
)

// @Summary Import ImageName
// @Produce  json
// @Param image formData file true "ImageName File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logger.Warn(err)
		appG.Response(http.StatusInternalServerError, define.ERROR, nil)
		return
	}

	if image == nil {
		appG.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusBadRequest, define.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logger.Warn(err)
		appG.Response(http.StatusInternalServerError, define.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		logger.Warn(err)
		appG.Response(http.StatusInternalServerError, define.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, define.SUCCESS, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
	})
}
