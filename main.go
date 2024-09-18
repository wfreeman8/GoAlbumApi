package main

import (
  "fmt"
  "net/http"
  "GoAlbumApi/controllers"
  "github.com/gin-gonic/gin"
)


var albumBasePath string = "albums"

func main() {
  router := gin.Default()
  router.LoadHTMLGlob("static/*")

  router.GET("/form", func(c *gin.Context) {
    c.HTML(http.StatusOK, "form.htm", gin.H{})
  })
  var albumsController = controllers.AlbumsController{albumBasePath}
  var albumController = controllers.AlbumController{albumBasePath}
  var albumImageController = controllers.AlbumImageController{albumBasePath}
  var albumImagesController = controllers.AlbumImagesController{albumBasePath}
  
  router.POST("/albums", albumsController.Post)
  router.GET("/albums", albumsController.Get)
  router.GET("/album/:albumPagename", albumController.Get)
  router.PUT("/album/:albumPagename", albumController.Put)
  router.POST("/album/:albumPagename/images", albumImagesController.Post)
  router.GET("/img/:albumPagename/:imgFileName", albumImageController.Get)
  fmt.Println("Starting Web Server")
  router.Run(":3000")

}