package controllers

import (
  "fmt"
  "reflect"
  "encoding/json"
  "net/http"
  "GoAlbumApi/models"
  "github.com/gin-gonic/gin"
)

type AlbumController struct {
  Config *models.Config
}

type PutAlbumFields struct {
  Title string `json:"title"`
  Author string `json:"author"`
  Description string `json:"description"`
  DateTaken string `json:"created_datetime"`
  Pagename string `json:"pagename"`
  FeaturedImage string `json:"featured_image"`
  ThumbnailSize string `json:"thumbnail_size"`
  LargeSize string `json:"large_size"`
  Resizes []models.ImageResize `json: "resizes"`
}

func (albumController *AlbumController) Get(ginContext *gin.Context) {
  albumPagename := ginContext.Param("albumPagename")

  album, err := models.FindAlbum(albumController.Config.AlbumBasePath, albumPagename)
  if err == nil {
    ginContext.JSON(http.StatusOK, album.GetAlbumFormatted())
    return
  }
  ginContext.JSON(http.StatusNotFound, gin.H{
    "message": "Album not found",
  })
}

func (albumController *AlbumController) Put(ginContext *gin.Context) {
  albumPagename := ginContext.Param("albumPagename")

  album, err := models.FindAlbum(albumController.Config.AlbumBasePath, albumPagename)
  if err == nil {
    albumIsDirty := false
    decoder := json.NewDecoder(ginContext.Request.Body)
    var updatedAlbumData PutAlbumFields
    err = decoder.Decode(&updatedAlbumData)
    if err == nil {
      putAlbumReflect := reflect.ValueOf(updatedAlbumData)
      putAlbumReflectType := putAlbumReflect.Type()

      albumReflect := reflect.ValueOf(&album)
      albumReflectElem := albumReflect.Elem()
      for i := 0; i < putAlbumReflect.NumField(); i++ {
        putAlbumField := putAlbumReflect.Field(i)
        fieldName := putAlbumReflectType.Field(i).Name
        value := putAlbumField.String()

        if putAlbumField.Kind() == reflect.String && value != "" {
          albumField := albumReflectElem.FieldByName(fieldName)
          if value != albumField.String() {
            albumField.SetString(value)
            albumIsDirty = true
          }
        } else if putAlbumField.Kind() == reflect.Slice {
          albumField := albumReflectElem.FieldByName(fieldName)
          fieldIsDirty := putAlbumField.Len() != albumField.Len() && putAlbumField.Len() > 0
          if putAlbumField.Len() == albumField.Len() && !fieldIsDirty {
            for i := 0; i < putAlbumField.Len(); i++ {
              if !putAlbumField.Index(i).Equal(albumField.Index(i)) {
                fieldIsDirty = true
              }
            }
          }

          if fieldName == "Resizes" && fieldIsDirty {
            album.DeleteResizeImages("")
            albumIsDirty = true
          }
        }
      }

      imagesWereChanged, imagesChangedErr := album.CheckAndResetImagesIndex()
      if imagesWereChanged && imagesChangedErr == nil {
        albumIsDirty = true
      }

      if albumIsDirty {
        err = album.Save()
      }

      if err == nil {
        ginContext.JSON(http.StatusOK, album.GetAlbumFormatted())
        return
      } else {
        fmt.Println(err)
      }
    }
    if err != nil {
      ginContext.JSON(http.StatusInternalServerError, gin.H{
        "success": false,
        "error": err.Error(),
      })
      return
    }
  }
  ginContext.JSON(http.StatusNotFound, gin.H{
    "message": "Album not found",
  })
}

