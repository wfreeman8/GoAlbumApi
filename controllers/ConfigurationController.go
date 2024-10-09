package controllers

import (
  "net/http"
  "encoding/json"
  "GoAlbumApi/models"
  "github.com/gin-gonic/gin"
)
type fn func()

type ConfigurationController struct {
  PostSetupFunc fn
  Config *models.Config
}

type ConfigPost struct {
  Password string `json:"password"`
  BasePath string `json:"albums_base_path"`
  GalleryTitle string `json:"gallery_title"`
}


func (configController *ConfigurationController) Get(ginContext *gin.Context) {
  config, err := models.GetConfiguration()
  if err != nil {
    ginContext.JSON(http.StatusNotFound, gin.H{
      "message": "Configuration Not Found",
    })
    return
  }
  ginContext.JSON(http.StatusOK, config.GetConfigFormatted())
}

func (configController *ConfigurationController) Post(ginContext *gin.Context) {
  postJsonDecoder := json.NewDecoder(ginContext.Request.Body)
  var postValues ConfigPost
  err := postJsonDecoder.Decode(&postValues)
  responseCode := http.StatusInternalServerError
  if err != nil {
    ginContext.JSON(responseCode, gin.H{
      "success": false,
    })
  }

  newConfiguration := models.Config{
    AlbumBasePath: postValues.BasePath,
    GalleryTitle: postValues.GalleryTitle,
  }
  newConfiguration.SetPassword(postValues.Password)
  if err = newConfiguration.UpsertBaseFolder(); err == nil {
    if err = newConfiguration.Save(true); err == nil {
      responseCode = http.StatusCreated
      configController.PostSetupFunc()
    }
  }
  ginContext.JSON(responseCode, gin.H{
    "success": responseCode == http.StatusCreated,
  })
}

func (configController *ConfigurationController) Put(ginContext *gin.Context) {
  postJsonDecoder := json.NewDecoder(ginContext.Request.Body)
  var postValues ConfigPost
  err := postJsonDecoder.Decode(&postValues)

  success := false
  if err != nil {
    ginContext.JSON(http.StatusOK, gin.H{
      "success": success,
    })
    return
  }
  if postValues.Password != "" {
    configController.Config.SetPassword(postValues.Password)
  }

  if postValues.GalleryTitle != "" {
    configController.Config.GalleryTitle = postValues.GalleryTitle
  }

  if err := configController.Config.Save(false); err == nil {
    success = true
  }
  ginContext.JSON(http.StatusOK, gin.H{
    "success": success,
  })
}