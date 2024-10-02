package controllers

import (
  "net/http"
  "GoAlbumApi/models"
  "github.com/gin-gonic/gin"
)

type AlbumImagesController struct {
  Config *models.Config
}

func (albumImagesController *AlbumImagesController) Post(ginContext *gin.Context) {
  albumPagename := ginContext.Param("albumPagename")
  album, err := models.FindAlbum(albumImagesController.Config.AlbumBasePath, albumPagename)
  if err == nil {
    imageFile, _ := ginContext.FormFile("new_image")
    err = album.SaveUploadedImage(imageFile)
    if err == nil {
      ginContext.JSON(http.StatusOK, album.GetAlbumFormatted().Images)
      return
    }
  }

  ginContext.JSON(http.StatusOK, gin.H{
    "success": false,
    "error": err.Error(),
  })
}
