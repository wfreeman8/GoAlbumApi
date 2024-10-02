package controllers

import (
  "net/http"
  "GoAlbumApi/models"
  "github.com/gin-gonic/gin"
  "regexp"
  "strconv"
)

type AlbumImageController struct{
  Config *models.Config
}

func (albumImageController *AlbumImageController) Get(ginContext *gin.Context) {
  albumPagename := ginContext.Param("albumPagename")
  album, err := models.FindAlbum(albumImageController.Config.AlbumBasePath, albumPagename)

  if err == nil {
    imgFileName := ginContext.Param("imgFileName")
    re := regexp.MustCompile(`([a-z0-9]+)-(\d+)x(\d+)(x1)?\.(jpg|png)`)
    imageFilenamePieces := re.FindAllStringSubmatch(imgFileName, -1)
    var realImageFilename string
    var imageResize models.ImageResize
    var resizeErr error

    if len(imageFilenamePieces) > 0 {
      realImageFilename = imageFilenamePieces[0][1] + "." + imageFilenamePieces[0][5]
      width, _ := strconv.Atoi(imageFilenamePieces[0][2])
      height, _ := strconv.Atoi(imageFilenamePieces[0][3])
      crop := false
      if imageFilenamePieces[0][4] == "x1" {
        crop = true
      }
      imageResize, resizeErr = album.FindResizeBySize(width, height, crop)
    } else {
      sizeNameRe := regexp.MustCompile(`([a-z0-9]+)-([A-Za-z0-9_]+)\.(jpg|png)`)
      imageFilenamePieces := sizeNameRe.FindAllStringSubmatch(imgFileName, -1)
      if len(imageFilenamePieces) > 0 {
        realImageFilename = imageFilenamePieces[0][1] + "." + imageFilenamePieces[0][3]
        imageResize, resizeErr = album.FindResizeByName(imageFilenamePieces[0][2])
      }
    }

    imageData, err := album.RetrieveImage(realImageFilename)
    if err == nil && resizeErr == nil {
      imgBytes := imageData.GetResizeBytes(imageResize)

      if err == nil {
        if len(imgBytes) == 0 {
          imgBytes = imageData.GetBytes()
        }
        ginContext.Data(http.StatusOK, "image/jpeg", imgBytes)
        return
      }
    }
  }
  ginContext.JSON(http.StatusNotFound, gin.H{
    "message": "Image Not Found",
  })
}

