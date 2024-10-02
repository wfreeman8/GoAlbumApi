package main

import (
  "fmt"
  "os"
  "net/http"
  "GoAlbumApi/controllers"
  "GoAlbumApi/models"
  "github.com/gin-gonic/gin"
)

var router = gin.Default()

// Temporary until setting up config.json for domains
func CORSMiddleware() gin.HandlerFunc {
  return func(ginContext *gin.Context) {
    ginContext.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    ginContext.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
    ginContext.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

    if ginContext.Request.Method == "OPTIONS" {
        ginContext.AbortWithStatus(204)
        return
    }

    ginContext.Next()
  }
}

// validate password for endpoints and abort if password is incorrect
func AuthorizationCheck(config *models.Config) gin.HandlerFunc {
  return func(ginContext *gin.Context) {
    if !config.ValidatePassword(ginContext.GetHeader("Authorization")) {
      ginContext.JSON(http.StatusUnauthorized, gin.H{
        "success": false,
        "error": "Invalid Authorization",
      })
      ginContext.Abort()
      return
    }

    ginContext.Next()
  }
}

func main() {
  router.LoadHTMLGlob("static/*.htm")
  currentDir, _ := os.Getwd()
  router.StaticFS("/js", http.Dir(currentDir + "/static/js"))

  router.GET("", func(c *gin.Context) {
    if models.HasConfiguration() {
      c.HTML(http.StatusOK, "form.htm", gin.H{})
    } else {
      c.HTML(http.StatusOK, "setup.htm", gin.H{})
    }
  })

  if models.HasConfiguration() {
    setupRESTRoutes()
    fmt.Println("Starting Web Server")
  } else {
    configurationController := controllers.ConfigurationController{PostSetupFunc: setupRESTRoutes}
    router.POST("/config", configurationController.Post)
  }
  router.Run(":3000")

}

func setupRESTRoutes() {
  config, _ := models.GetConfiguration()
  albumsController := controllers.AlbumsController{config}
  albumController := controllers.AlbumController{config}
  albumImageController := controllers.AlbumImageController{config}
  albumImagesController := controllers.AlbumImagesController{config}
  configurationController := controllers.ConfigurationController{
    Config: config,
  }
  router.Use(CORSMiddleware())
  router.GET("/config", configurationController.Get)
  router.GET("/albums", albumsController.Get)
  router.GET("/album/:albumPagename", albumController.Get)
  router.GET("/img/:albumPagename/:imgFileName", albumImageController.Get)

  authRequiredRoutes := router.Group("/")
  authRequiredRoutes.Use(AuthorizationCheck(config))
  {
    authRequiredRoutes.PUT("/config", configurationController.Put)
    authRequiredRoutes.POST("/albums", albumsController.Post)
    authRequiredRoutes.PUT("/album/:albumPagename", albumController.Put)
    authRequiredRoutes.POST("/album/:albumPagename/images", albumImagesController.Post)
  }

}