package controllers

import (
  "fmt"
  "os"
  "io"
  "errors"
  "net/http"
  "encoding/json"
  "github.com/google/uuid"
  "GoAlbumApi/models"
  "github.com/gin-gonic/gin"
)

type AlbumSummaryResponse struct {
  AlbumTitle string `json:"title"`
  Pagename string `json:"pagename"`
  Description string `json:"description"`
  ImageCount int `json:"imagecount"`
  Author string `json:"author"`
}

type AlbumsController struct {
  Config *models.Config
}

func (albumsController *AlbumsController) Get(ginContext *gin.Context) {
  albumsEntries, err := os.ReadDir(albumsController.Config.AlbumBasePath)
  var albumResponseCollection = []AlbumSummaryResponse{}

  if err != nil {
    ginContext.JSON(http.StatusBadRequest, gin.H{
      "success": false,
      "error": err.Error(),
    })
    return
  }

  for _, albumEntry := range albumsEntries {
    albumDirPath := albumsController.Config.AlbumBasePath + "/" + albumEntry.Name()
    albumFolderMeta, err := os.Stat(albumDirPath)
    if err != nil || !albumFolderMeta.IsDir(){
      continue
    }

    albumJsonPath := albumDirPath + "/album.json"
    albumJsonFile, err := os.Open(albumJsonPath)
    defer albumJsonFile.Close()
    if err != nil {
      fmt.Println(err)
      continue
    }

    albumJsonContent, _ := io.ReadAll(albumJsonFile)
    var albumData models.Album
    json.Unmarshal(albumJsonContent, &albumData)

    AlbumSummaryResponse := AlbumSummaryResponse{
      AlbumTitle: albumData.Title,
      Pagename: albumData.Pagename,
      Description: albumData.Description,
      ImageCount: len(albumData.Images),
      Author: albumData.Author,
    }
    albumResponseCollection = append(albumResponseCollection, AlbumSummaryResponse)
  }
  ginContext.JSON(http.StatusOK, albumResponseCollection)
}

func (albumsController *AlbumsController) Post(ginContext *gin.Context) {
  decoder := json.NewDecoder(ginContext.Request.Body)
  var album models.Album
  var err error
  err = decoder.Decode(&album)
  if err != nil {
    ginContext.JSON(http.StatusBadRequest, gin.H{
      "success": false,
      "error": err.Error(),
    })
    return
  }

  albumPath := albumsController.Config.AlbumBasePath + "/" + album.Pagename
  albumJsonPath := albumPath + "/album.json"
  fmt.Println(albumPath)

  albumInfo, err := os.Stat(albumPath)

  if errors.Is(err, os.ErrNotExist) {
    err = os.Mkdir(albumPath, 0754)
    if err == nil {
      albumInfo, err = os.Stat(albumPath)
    }
  }
  if err != nil {
    ginContext.JSON(http.StatusBadRequest, gin.H{
      "success": false,
      "error": err.Error(),
    })
    return
  }

  if albumInfo.IsDir() {
    album.SetAlbumPath(albumPath)
    if _, err = os.Stat(albumJsonPath); err == nil {
      ginContext.JSON(http.StatusConflict, gin.H{
        "success": false,
        "error": "Album already exist please use PUT to update",
      })
      return
    }

    album.AlbumId = (uuid.New()).String()
    album.IndexImagesWithDirectory()
    album.Save()

    ginContext.JSON(http.StatusCreated, album)
    return
  }
}
